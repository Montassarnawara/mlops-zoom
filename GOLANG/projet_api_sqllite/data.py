import pandas as pd
import sqlite3

prod = pd.read_csv("product.csv")
conn = sqlite3.connect("product.db")
prod.to_sql("product", conn, if_exists="replace", index=False)
conn.close()