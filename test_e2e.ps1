# test_e2e.ps1 - 工单系统端到端测试
$BaseUrl = "http://localhost:8000/api/v1"
$ApiKey = "sk-test-key-12345"
$Headers = @{ "X-API-Key" = $ApiKey; "Content-Type" = "application/json" }

$pass = 0; $fail = 0
function Test-Case {
    param($Name, $Expected, $Actual, $Body)
    $script:total++
    if ($Expected -eq $Actual) {
        Write-Host "  ✅ $Name" -ForegroundColor Green
        $script:pass++
    } else {
        Write-Host "  ❌ $Name (expected=$Expected, got=$Actual)" -ForegroundColor Red
        if ($Body) { Write-Host "     $Body" -ForegroundColor Gray }
        $script:fail++
    }
}

function Invoke-Json {
    param($Method, $Path, $Body = $null, $Headers = $null)
    $params = @{
        Uri = $BaseUrl + $Path
        Method = $Method
        UseBasicParsing = $true
        TimeoutSec = 10
    }
    if ($Headers) { $params.Headers = $Headers }
    if ($Body) { $params.Body = ($Body | ConvertTo-Json -Depth 10) }
    try {
        $r = Invoke-WebRequest @params
        return @{ Status = [int]$r.StatusCode; Body = ($r.Content | ConvertFrom-Json) }
    } catch {
        $code = 0
        if ($_.Exception.Response) { $code = [int]$_.Exception.Response.StatusCode }
        $body = $null
        try { $body = ($_.Exception.Response.GetResponseStream() | New-Object IO.StreamReader).ReadToEnd() | ConvertFrom-Json } catch {}
        return @{ Status = $code; Body = $body; Error = $_.Exception.Message }
    }
}

Write-Host "`n=== 1. 健康检查 ===" -ForegroundColor Cyan
$h = Invoke-WebRequest "http://localhost:8000/health" -UseBasicParsing
Test-Case "GET /health" 200 ([int]$h.StatusCode)

Write-Host "`n=== 2. 鉴权测试 ===" -ForegroundColor Cyan
$r = Invoke-Json POST "/tickets" @{ question = "x"; user_id = "y" } @{ "Content-Type" = "application/json" }
Test-Case "无 API Key → 401" 401 $r.Status

$r = Invoke-Json POST "/tickets" @{ question = "x"; user_id = "y" } @{ "Content-Type" = "application/json"; "X-API-Key" = "wrong-key" }
Test-Case "错误 API Key → 401" 401 $r.Status

Write-Host "`n=== 3. 提交工单 ===" -ForegroundColor Cyan
$body1 = @{
    question = "石河子大学2026年计算机学院宿舍分配在哪？"
    user_id = "sess_test_001"
    source = "hiagent_chat"
    rag_result = ""
}
$r = Invoke-Json POST "/tickets" $body1 $Headers
Test-Case "POST /tickets → 200" 200 $r.Status
$ticketId = $r.Body.ticket_id
Write-Host "     工单号: $ticketId" -ForegroundColor Gray
Test-Case "返回 ticket_id 非空" $true ($ticketId -match "^T\d{8}-[A-Z0-9]{6}$")
Test-Case "status = created" "created" $r.Body.status
Test-Case "message 包含工单号" $true ($r.Body.message -match $ticketId)

Write-Host "`n=== 4. 5 分钟内重复提交（同 user+question）===" -ForegroundColor Cyan
$r = Invoke-Json POST "/tickets" $body1 $Headers
Test-Case "重复提交 → 200" 200 $r.Status
Test-Case "返回相同工单号" $ticketId $r.Body.ticket_id
Test-Case "message 提示已记录" $true ($r.Body.message -match "已经记录过")

Write-Host "`n=== 5. 查询工单详情 ===" -ForegroundColor Cyan
$r = Invoke-Json GET "/tickets/$ticketId" -Headers $Headers
Test-Case "GET /tickets/{id} → 200" 200 $r.Status
Test-Case "status = pending" "pending" $r.Body.status
Test-Case "question 匹配" "石河子大学2026年计算机学院宿舍分配在哪？" $r.Body.question
Test-Case "answer 字段为 null" $null $r.Body.answer
Test-Case "answered_at 字段为 null" $null $r.Body.answered_at

Write-Host "`n=== 6. 人工答复 ===" -ForegroundColor Cyan
$ansBody = @{
    answer = "计算机学院2026级新生安排在北区3号公寓，6人间，上床下桌，有独立卫生间。"
    operator = "张管理员"
    sync_to_kb = $true
}
$r = Invoke-Json POST "/tickets/$ticketId/answer" $ansBody $Headers
Test-Case "POST /tickets/{id}/answer → 200" 200 $r.Status
Test-Case "success = true" $true $r.Body.success
Test-Case "message 包含工单号" $true ($r.Body.message -match $ticketId)

# 验证答复内容
$r = Invoke-Json GET "/tickets/$ticketId" -Headers $Headers
Test-Case "查询后 status = answered" "answered" $r.Body.status
Test-Case "answer 字段已填充" $true ($r.Body.answer.Length -gt 0)
Test-Case "answered_by = 张管理员" "张管理员" $r.Body.answered_by
Test-Case "answered_at 已填充" $true ($null -ne $r.Body.answered_at)

Write-Host "`n=== 7. 关闭工单 ===" -ForegroundColor Cyan
$closeBody = @{ reason = "测试完成，主动关闭" }
$r = Invoke-Json POST "/tickets/$ticketId/close" $closeBody $Headers
Test-Case "POST /tickets/{id}/close → 200" 200 $r.Status

$r = Invoke-Json GET "/tickets/$ticketId" -Headers $Headers
Test-Case "关闭后 status = closed" "closed" $r.Body.status
Test-Case "close_reason 已记录" "测试完成，主动关闭" $r.Body.close_reason

# 关闭后再次答复应失败
$r = Invoke-Json POST "/tickets/$ticketId/answer" $ansBody $Headers
Test-Case "关闭后答复 → 400" 400 $r.Status

Write-Host "`n=== 8. 查询不存在的工单 ===" -ForegroundColor Cyan
$r = Invoke-Json GET "/tickets/T99999999-XXXXXX" -Headers $Headers
Test-Case "GET 不存在工单 → 404" 404 $r.Status

Write-Host "`n=== 9. 工单列表 ===" -ForegroundColor Cyan
$r = Invoke-Json GET "/tickets?page=1&page_size=10" -Headers $Headers
Test-Case "GET /tickets → 200" 200 $r.Status
Test-Case "total >= 1" $true ($r.Body.total -ge 1)
Test-Case "items 数组非空" $true ($r.Body.items.Count -ge 1)

$r = Invoke-Json GET "/tickets?status=closed" -Headers $Headers
Test-Case "按状态筛选 closed" $true ($r.Body.items | Where-Object { $_.status -eq "closed" } | Measure-Object).Count -ge 1

Write-Host "`n=== 10. 参数校验 ===" -ForegroundColor Cyan
$r = Invoke-Json POST "/tickets" @{ user_id = "y" } $Headers
Test-Case "缺少 question → 422" 422 $r.Status

$r = Invoke-Json POST "/tickets" @{ question = ""; user_id = "y" } $Headers
Test-Case "空字符串 question → 422" 422 $r.Status

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host " 通过: $pass / $total" -ForegroundColor Green
Write-Host " 失败: $fail / $total" -ForegroundColor Red
Write-Host "========================================`n" -ForegroundColor Cyan

if ($fail -gt 0) { exit 1 } else { exit 0 }
