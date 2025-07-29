import pandas as pd
import psycopg2

# Charger les données
df = pd.read_csv("college.csv")

# Connexion à PostgreSQL
conn = psycopg2.connect(
    host="localhost",
    dbname="ma_base",
    user="montassar",
    password="123mont@456"
)
cur = conn.cursor()

# Création de la table en snake_case
cur.execute("""
    CREATE TABLE IF NOT EXISTS college (
        college_id TEXT PRIMARY KEY,
        iq INTEGER,
        prev_sem_result FLOAT,
        cgpa FLOAT,
        academic_performance INTEGER,
        internship_experience TEXT,
        extra_curricular_score INTEGER,
        communication_skills INTEGER,
        projects_completed INTEGER,
        placement TEXT
    );
""")
conn.commit()

# Insertion ligne par ligne
for _, row in df.iterrows():
    cur.execute("""
        INSERT INTO college (
            college_id, iq, prev_sem_result, cgpa, academic_performance,
            internship_experience, extra_curricular_score, communication_skills,
            projects_completed, placement
        )
        VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
        ON CONFLICT (college_id) DO NOTHING;
    """, tuple(row))

conn.commit()
cur.close()
conn.close()

print("✅ Données insérées avec succès dans la table 'college'.")
