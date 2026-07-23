"""
debug_test.py - 单步调试，复现 test_e2e 的 step 6 问题
"""
import urllib.request, urllib.error, json

BASE = "http://localhost:8000/api/v1"


def call(method, path, body=None, headers=None, base=BASE):
    url = base + path
    data = json.dumps(body).encode("utf-8") if body is not None else None
    h = headers if headers is not None else {}
    # debug print
    print(f"  -> {method} {path}")
    print(f"     headers: {h}")
    req = urllib.request.Request(url, data=data, method=method, headers=h)
    try:
        with urllib.request.urlopen(req, timeout=5) as r:
            return r.status, json.loads(r.read().decode("utf-8"))
    except urllib.error.HTTPError as e:
        try:
            return e.code, json.loads(e.read().decode("utf-8"))
        except:
            return e.code, None


# Step 3: admin login
print("=== Step 3: admin login ===")
code, resp = call("POST", "/auth/login", {"username": "admin", "password": "admin123"},
                  {"Content-Type": "application/json"})
print(f"  result: {code}, token: {bool(resp and resp.get('token'))}")
admin_token = resp.get("token")
ADMIN_HDR = {"Authorization": f"Bearer {admin_token}", "Content-Type": "application/json"}
print(f"  ADMIN_HDR set: {ADMIN_HDR}")

# Step 4: submit ticket
print("\n=== Step 4: submit ticket ===")
code, resp = call("POST", "/tickets", {"question": f"debug test", "user_id": "u1"},
                  {"X-API-Key": "sk-test-key-12345", "Content-Type": "application/json"})
print(f"  result: {code}, ticket_id: {resp and resp.get('ticket_id')}")
ticket_id = resp.get("ticket_id")

# Step 6: get ticket with ADMIN_HDR
print("\n=== Step 6: get ticket with ADMIN_HDR ===")
print(f"  ADMIN_HDR is: {ADMIN_HDR}")
code, resp = call("GET", f"/tickets/{ticket_id}", ADMIN_HDR)
print(f"  result: {code}, resp: {resp}")
