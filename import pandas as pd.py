import pandas as pd
from sqlalchemy import create_engine

# Charger le CSV
college = pd.read_csv("college_student_placement_dataset.csv")

# Connexion à PostgreSQL
engine = create_engine("postgresql://postgres:motdepasse@localhost:5432/ma_base")

# Envoi à PostgreSQL
college.to_sql("personnes", engine, if_exists="replace", index=False)

print("Données insérées avec succès dans PostgreSQL")
