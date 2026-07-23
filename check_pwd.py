import sqlite3
import bcrypt
db = r"C:\Users\A\Desktop\迎新智能体\智能体插件\ticket-system\backend\tickets.db"
c = sqlite3.connect(db)
row = c.execute("SELECT username, password_hash FROM users WHERE username='admin'").fetchone()
print("row:", row)
print("admin123 verify:", bcrypt.checkpw(b"admin123", row[1].encode()))
print("newAdminPass123 verify:", bcrypt.checkpw(b"newAdminPass123", row[1].encode()))
