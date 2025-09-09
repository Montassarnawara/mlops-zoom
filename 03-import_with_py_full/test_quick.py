"""
🧪 TEST RAPIDE - Vérification du script de base de données
============================================================

Ce script teste la création de quelques tables pour valider la logique
avant d'exécuter le script complet.
"""

import psycopg2
from psycopg2.extensions import ISOLATION_LEVEL_AUTOCOMMIT
from faker import Faker
import random
from datetime import timedelta, date

fake = Faker()

def test_connection():
    """Tester la connexion PostgreSQL"""
    try:
        conn = psycopg2.connect(
            dbname="postgres",
            user="myuser",
            password="mypassword",
            host="localhost",
            port="5433"
        )
        conn.close()
        print("✅ Connexion PostgreSQL réussie")
        return True
    except Exception as e:
        print(f"❌ Erreur de connexion: {e}")
        return False

def test_table_creation():
    """Tester la création de quelques tables"""
    if not test_connection():
        return False
    
    try:
        # Créer une base de test
        conn = psycopg2.connect(
            dbname="postgres",
            user="myuser",
            password="mypassword",
            host="localhost",
            port="5433"
        )
        conn.set_isolation_level(ISOLATION_LEVEL_AUTOCOMMIT)
        cur = conn.cursor()
        
        # Supprimer la base de test si elle existe
        cur.execute("DROP DATABASE IF EXISTS test_hotel_db")
        cur.execute("CREATE DATABASE test_hotel_db")
        print("✅ Base de test créée")
        
        cur.close()
        conn.close()
        
        # Se connecter à la base de test
        conn = psycopg2.connect(
            dbname="test_hotel_db",
            user="myuser",
            password="mypassword",
            host="localhost",
            port="5433"
        )
        cur = conn.cursor()
        
        # Créer quelques tables de test
        cur.execute("""
            CREATE TABLE COUNTRY (
                id SERIAL PRIMARY KEY,
                name VARCHAR(100) NOT NULL,
                code VARCHAR(10) UNIQUE NOT NULL
            );
            
            CREATE TABLE CITY (
                id SERIAL PRIMARY KEY,
                citycode VARCHAR(20) NOT NULL,
                name VARCHAR(100) NOT NULL,
                countryname VARCHAR(100),
                countryid INT REFERENCES COUNTRY(id) ON DELETE CASCADE
            );
            
            CREATE TABLE HOTEL (
                hotel_key VARCHAR(20) PRIMARY KEY,
                name VARCHAR(100) NOT NULL,
                city VARCHAR(100),
                country VARCHAR(100),
                stars INT CHECK (stars BETWEEN 1 AND 5)
            );
        """)
        
        print("✅ Tables de test créées")
        
        # Insérer quelques données
        cur.execute("INSERT INTO COUNTRY (name, code) VALUES ('France', 'FR')")
        cur.execute("INSERT INTO COUNTRY (name, code) VALUES ('Espagne', 'ES')")
        
        cur.execute("INSERT INTO CITY (citycode, name, countryname, countryid) VALUES ('PAR', 'Paris', 'France', 1)")
        cur.execute("INSERT INTO CITY (citycode, name, countryname, countryid) VALUES ('MAD', 'Madrid', 'Espagne', 2)")
        
        cur.execute("INSERT INTO HOTEL (hotel_key, name, city, country, stars) VALUES ('HOT001', 'Hotel Test Paris', 'Paris', 'France', 4)")
        cur.execute("INSERT INTO HOTEL (hotel_key, name, city, country, stars) VALUES ('HOT002', 'Hotel Test Madrid', 'Madrid', 'Espagne', 5)")
        
        conn.commit()
        print("✅ Données de test insérées")
        
        # Vérifier les données
        cur.execute("""
            SELECT h.hotel_key, h.name, h.city, h.stars 
            FROM HOTEL h 
            ORDER BY h.hotel_key
        """)
        
        hotels = cur.fetchall()
        print("\n📋 Hôtels de test créés:")
        for hotel in hotels:
            print(f"   {hotel[0]} - {hotel[1]} ({hotel[2]}, {hotel[3]} étoiles)")
        
        # Nettoyer
        cur.close()
        conn.close()
        
        # Supprimer la base de test
        conn = psycopg2.connect(
            dbname="postgres",
            user="myuser",
            password="mypassword",
            host="localhost",
            port="5433"
        )
        conn.set_isolation_level(ISOLATION_LEVEL_AUTOCOMMIT)
        cur = conn.cursor()
        cur.execute("DROP DATABASE test_hotel_db")
        cur.close()
        conn.close()
        
        print("✅ Base de test nettoyée")
        return True
        
    except Exception as e:
        print(f"❌ Erreur lors du test: {e}")
        return False

def main():
    print("🧪 TEST RAPIDE DU SYSTÈME DE BASE DE DONNÉES")
    print("=" * 50)
    
    if test_table_creation():
        print("\n🎉 TOUS LES TESTS RÉUSSIS !")
        print("=" * 50)
        print("✅ Connexion PostgreSQL fonctionnelle")
        print("✅ Création de tables réussie")
        print("✅ Insertion de données réussie")
        print("✅ Contraintes de clés étrangères respectées")
        print("\n🚀 Le script complet devrait fonctionner correctement.")
        print("   Exécutez: python3 complete_hotel_db.py")
    else:
        print("\n❌ ÉCHEC DES TESTS")
        print("Vérifiez votre configuration PostgreSQL")

if __name__ == "__main__":
    main()
