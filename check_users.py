import sqlite3
db = r"C:\Users\A\Desktop\迎新智能体\智能体插件\ticket-system\backend\tickets.db"
c = sqlite3.connect(db)
print("All users:")
for u in c.execute("SELECT id, username, role, active, must_change_pwd, created_at, length(password_hash) FROM users"):
    print(" ", u)
