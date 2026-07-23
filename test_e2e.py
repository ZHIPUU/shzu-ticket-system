"""
test_e2e.py - 工单系统全量端到端测试

覆盖：
  1. 工单 CRUD（37 用例）
  2. JWT 鉴权 + 用户管理（46 用例）

运行前置：
  - 后端已启动
  - 默认 admin/admin123 存在（首次启动自动创建）
  - 干净状态：trash tickets.db 再重启
"""
import json
import time
import urllib.request
import urllib.error

BASE = "http://localhost:8000/api/v1"
HEALTH = "http://localhost:8000/health"
RUN_ID = str(int(time.time()))

pass_cnt = 0
fail_cnt = 0
total_cnt = 0


def call(method, path, body=None, headers=None, base=BASE):
    url = base + path
    data = json.dumps(body).encode("utf-8") if body is not None else None
    h = headers if headers is not None else {}
    req = urllib.request.Request(url, data=data, method=method, headers=h)
    try:
        with urllib.request.urlopen(req, timeout=5) as r:
            return r.status, json.loads(r.read().decode("utf-8"))
    except urllib.error.HTTPError as e:
        try:
            body_resp = json.loads(e.read().decode("utf-8"))
        except Exception:
            body_resp = None
        return e.code, body_resp
    except Exception as e:
        return 0, str(e)


def check(name, expected, actual, extra=""):
    global pass_cnt, fail_cnt, total_cnt
    total_cnt += 1
    if expected == actual:
        pass_cnt += 1
        print(f"  PASS  {name}")
    else:
        fail_cnt += 1
        print(f"  FAIL  {name}  (expected={expected!r}, got={actual!r}) {extra}")


def check_true(name, condition, extra=""):
    global pass_cnt, fail_cnt, total_cnt
    total_cnt += 1
    if condition:
        pass_cnt += 1
        print(f"  PASS  {name}")
    else:
        fail_cnt += 1
        print(f"  FAIL  {name}  {extra}")


def section(t):
    print(f"\n=== {t} ===")


# ════════════════════════════════════════════════════════════════
# 第一部分：工单 CRUD（含双轨鉴权）
# ════════════════════════════════════════════════════════════════
print("\n" + "=" * 60)
print("  PART 1: 工单 API（含鉴权）")
print("=" * 60)

# ── 1. 健康检查 ──
section("1. Health Check")
try:
    with urllib.request.urlopen(HEALTH, timeout=3) as r:
        check("GET /health -> 200", 200, r.status)
except Exception:
    check("server up", 200, 0)

# ── 2. 鉴权测试 ──
section("2. Auth")
code, _ = call("POST", "/tickets", {"question": "x", "user_id": "y"},
               {"Content-Type": "application/json"})
check("no API key + no JWT -> 401", 401, code)
code, _ = call("POST", "/tickets", {"question": "x", "user_id": "y"},
               {"Content-Type": "application/json", "X-API-Key": "wrong"})
check("wrong API key -> 401", 401, code)

APIKEY_HDR = {"X-API-Key": "sk-test-key-12345", "Content-Type": "application/json"}

# ── 3. admin 登录拿 JWT ──
section("3. Admin Login")
code, resp = call("POST", "/auth/login", {"username": "admin", "password": "admin123"})
if code != 200:
    print(f"❌ admin login failed ({code}); 请先 trash tickets.db 并重启后端")
    print(resp)
    exit(1)
admin_token = resp.get("token")
ADMIN_HDR = {"Authorization": f"Bearer {admin_token}", "Content-Type": "application/json"}
check("admin login -> 200", 200, code)
check_true("token returned", bool(admin_token))
check("user.role = admin", "admin", resp["user"]["role"])

# ── 4. 提交工单（API Key 给智能体）──
section("4. Submit Ticket (API Key)")
q1 = f"石河子大学2026年计算机学院宿舍分配在哪？[{RUN_ID}]"
body1 = {"question": q1, "user_id": f"sess_{RUN_ID}", "source": "hiagent_chat", "rag_result": ""}
code, resp = call("POST", "/tickets", body1, APIKEY_HDR)
check("POST /tickets (API key) -> 200", 200, code)
ticket_id = resp.get("ticket_id")
check_true("ticket_id format", ticket_id and len(ticket_id) == 16 and ticket_id.startswith("T"))

# ── 5. 5 分钟内去重 ──
section("5. Dedup")
code, resp = call("POST", "/tickets", body1, APIKEY_HDR)
check("dup submit -> 200", 200, code)
check("same ticket_id", ticket_id, resp.get("ticket_id"))

# ── 6. 查询详情 ──
section("6. Get Ticket")
code, resp = call("GET", f"/tickets/{ticket_id}", headers=ADMIN_HDR)
check("GET /tickets/{id} (JWT) -> 200", 200, code)
check("status = pending", "pending", resp.get("status"))
check("answer = null", None, resp.get("answer"))

# ── 7. 答复工单（JWT 用户，operator 由系统填）──
section("7. Answer Ticket (JWT user)")
ans = {"answer": "计算机学院2026级新生在北区3号公寓，6人间。", "sync_to_kb": True}
code, resp = call("POST", f"/tickets/{ticket_id}/answer", ans, ADMIN_HDR)
check("POST /answer -> 200", 200, code)
check_true("success = true", resp.get("success") is True)

# 验证答复人来自 JWT 用户
code, resp = call("GET", f"/tickets/{ticket_id}", headers=ADMIN_HDR)
check("status = answered", "answered", resp.get("status"))
check_true("answer populated", (resp.get("answer") or "").startswith("计算机学院"))
check("answered_by = 系统管理员", "系统管理员", resp.get("answered_by"))

# ── 8. 关闭工单 ──
section("8. Close Ticket")
code, _ = call("POST", f"/tickets/{ticket_id}/close", {"reason": "测试完成"}, ADMIN_HDR)
check("POST /close -> 200", 200, code)
code, resp = call("GET", f"/tickets/{ticket_id}", headers=ADMIN_HDR)
check("after close status", "closed", resp.get("status"))

# ── 9. 404 + 列表 ──
section("9. NotFound & List")
code, _ = call("GET", "/tickets/T99999999-XXXXXX", headers=ADMIN_HDR)
check("nonexistent -> 404", 404, code)
code, resp = call("GET", "/tickets?page=1&page_size=10", headers=ADMIN_HDR)
check("GET /tickets -> 200", 200, code)
check_true("total >= 1", resp.get("total", 0) >= 1)

# ── 10. 参数校验 ──
section("10. Validation")
code, _ = call("POST", "/tickets", {"user_id": "y"}, APIKEY_HDR)
check("missing question -> 422", 422, code)

# ── 11. 多源渠道 ──
section("11. Source Channels")
for src in ["wechat_service", "feishu", "yiban"]:
    code, _ = call("POST", "/tickets", {
        "question": f"测试-{src}[{RUN_ID}]", "user_id": f"sess_{src}_{RUN_ID}",
        "source": src, "rag_result": "无结果",
    }, APIKEY_HDR)
    check(f"source={src} -> 200", 200, code)
code, _ = call("POST", "/tickets", {
    "question": f"bad[{RUN_ID}]", "user_id": "x", "source": "fake",
}, APIKEY_HDR)
check("invalid source -> 422", 422, code)


# ════════════════════════════════════════════════════════════════
# 第二部分：JWT 鉴权 + 用户管理
# ════════════════════════════════════════════════════════════════
print("\n" + "=" * 60)
print("  PART 2: JWT 鉴权 + 用户管理")
print("=" * 60)

# ── 12. 错误密码 ──
section("12. Wrong Password")
code, resp = call("POST", "/auth/login", {"username": "admin", "password": "wrong-pwd"})
check("wrong pwd -> 401", 401, code)
check("error_code", "INVALID_CREDENTIALS", (resp or {}).get("error_code"))

# ── 13. /auth/me ──
section("13. /auth/me")
code, resp = call("GET", "/auth/me", headers=ADMIN_HDR)
check("me -> 200", 200, code)
check("me.username = admin", "admin", (resp or {}).get("user", {}).get("username"))
code, _ = call("GET", "/auth/me")
check("me no token -> 401", 401, code)

# ── 14. 弱密码校验 ──
section("14. Weak Password Rejected")
code, _ = call("POST", "/auth/change-password",
               {"old_password": "admin123", "new_password": "short"},
               ADMIN_HDR)
check("short pwd -> 422", 422, code)
code, _ = call("POST", "/auth/change-password",
               {"old_password": "admin123", "new_password": "nodigitshere"},
               ADMIN_HDR)
check("no digits pwd -> 400", 400, code)

# ── 15. 改密码再登录 ──
section("15. Change Password")
new_pwd = "newAdminPass123"
code, _ = call("POST", "/auth/change-password",
               {"old_password": "admin123", "new_password": new_pwd}, ADMIN_HDR)
check("change pwd -> 200", 200, code)
code, _ = call("POST", "/auth/login", {"username": "admin", "password": "admin123"})
check("old pwd login -> 401", 401, code)
code, resp = call("POST", "/auth/login", {"username": "admin", "password": new_pwd})
check("new pwd login -> 200", 200, code)
admin_token = resp["token"]
ADMIN_HDR = {"Authorization": f"Bearer {admin_token}", "Content-Type": "application/json"}

# ── 16. 新建 staff 用户 ──
section("16. Create Staff")
staff_user = f"staff_{RUN_ID}"
code, resp = call("POST", "/users",
                  {"username": staff_user, "password": "staff123", "role": "staff",
                   "display_name": "测试小王", "email": f"{staff_user}@shzu.edu.cn"},
                  ADMIN_HDR)
check("create staff -> 201", 201, code)
staff_id = (resp or {}).get("user", {}).get("id")
code, _ = call("POST", "/users",
               {"username": staff_user, "password": "staff123", "role": "staff"}, ADMIN_HDR)
check("duplicate username -> 409", 409, code)

# staff 首登 + 改密码
code, resp = call("POST", "/auth/login",
                  {"username": staff_user, "password": "staff123"})
check("staff login -> 200", 200, code)
check_true("staff must_change_pwd = true", resp.get("must_change_password") is True)
staff_token = resp["token"]
STAFF_HDR = {"Authorization": f"Bearer {staff_token}", "Content-Type": "application/json"}
code, _ = call("POST", "/auth/change-password",
               {"old_password": "staff123", "new_password": "staffNew123"}, STAFF_HDR)
check("staff change pwd -> 200", 200, code)
code, resp = call("POST", "/auth/login",
                  {"username": staff_user, "password": "staffNew123"})
staff_token = resp["token"]
STAFF_HDR = {"Authorization": f"Bearer {staff_token}", "Content-Type": "application/json"}

# ── 17. 权限：staff 不能管用户 ──
section("17. RBAC")
code, _ = call("GET", "/users", headers=STAFF_HDR)
check("staff GET /users -> 403", 403, code)

# ── 18. admin 列表/编辑/禁用用户 ──
section("18. Admin: List/Edit/Disable")
code, resp = call("GET", "/users", headers=ADMIN_HDR)
check("admin list users -> 200", 200, code)
check_true("total >= 2", resp.get("total", 0) >= 2)

code, _ = call("PATCH", f"/users/{staff_id}", {"active": False}, ADMIN_HDR)
check("disable user -> 200", 200, code)
code, _ = call("POST", "/auth/login",
               {"username": staff_user, "password": "staffNew123"})
check("disabled user login -> 403", 403, code)
code, _ = call("PATCH", f"/users/{staff_id}", {"active": True}, ADMIN_HDR)
check("re-enable user -> 200", 200, code)

# 防自禁用
code, resp = call("PATCH", "/users/1", {"active": False}, ADMIN_HDR)
check("self-disable -> 400", 400, code)
check_true("error_code = SELF_DISABLE", (resp or {}).get("error_code") == "SELF_DISABLE")

# 防降级最后 admin
code, resp = call("PATCH", "/users/1", {"role": "staff"}, ADMIN_HDR)
check("demote last admin -> 400", 400, code)
check_true("error_code = LAST_ADMIN", (resp or {}).get("error_code") == "LAST_ADMIN")


# ─── 总结 ───
print(f"\n{'=' * 60}")
print(f"  PASS: {pass_cnt} / {total_cnt}")
print(f"  FAIL: {fail_cnt} / {total_cnt}")
print(f"{'=' * 60}\n")
exit(0 if fail_cnt == 0 else 1)
