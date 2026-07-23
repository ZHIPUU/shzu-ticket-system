import sqlite3
import sys
db = r"C:\Users\A\Desktop\迎新智能体\智能体插件\ticket-system\backend\tickets.db"
c = sqlite3.connect(db)
print("users:")
for u in c.execute("SELECT id, username, role, active, length(password_hash) FROM users"):
    print(" ", u)
print("tickets count:", c.execute("SELECT count(*) FROM tickets").fetchone()[0])
