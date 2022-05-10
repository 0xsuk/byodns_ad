#script to copy /etc/pihole/gravity.db gravity to domains.db gravity
import sqlite3
con = sqlite3.connect("/etc/pihole/gravity.db")
cur = con.cursor()

noc = sqlite3.connect("domains.db")
ruc = noc.cursor()

i = 0
for row in cur.execute("select domain from gravity"):
    i = i+1
    ruc.execute("insert into gravity values (?,?,?)", row)


print(i, "domains written")
