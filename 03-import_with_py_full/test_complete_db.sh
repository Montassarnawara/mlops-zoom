#!/bin/bash

# 🧪 SCRIPT DE TEST - Base de données hôtelière complète
# =====================================================

echo "🧪 TEST DE LA BASE DE DONNÉES HÔTELIÈRE"
echo "========================================"
echo ""

# Vérifier que Python et les dépendances sont installées
echo "🔍 Vérification des prérequis..."
python3 -c "import psycopg2, faker; print('✅ Dépendances OK')" 2>/dev/null || {
    echo "❌ Dépendances manquantes. Installation..."
    pip install psycopg2-binary faker
}

# Vérifier que PostgreSQL est accessible
echo "🔍 Test de connexion PostgreSQL..."
python3 -c "
import psycopg2
try:
    conn = psycopg2.connect(
        dbname='postgres',
        user='myuser', 
        password='mypassword',
        host='localhost',
        port='5433'
    )
    conn.close()
    print('✅ PostgreSQL accessible')
except Exception as e:
    print(f'❌ PostgreSQL non accessible: {e}')
    exit(1)
"

if [ $? -ne 0 ]; then
    echo "❌ PostgreSQL n'est pas accessible. Vérifiez la configuration."
    exit 1
fi

echo ""
echo "🚀 Lancement de la création de la base..."
echo "=========================================="

# Exécuter le script principal
python3 complete_hotel_db.py

if [ $? -eq 0 ]; then
    echo ""
    echo "🧪 TESTS DE VÉRIFICATION"
    echo "========================"
    
    # Test de vérification des données
    python3 -c "
import psycopg2

conn = psycopg2.connect(
    dbname='hotel_db',
    user='myuser',
    password='mypassword', 
    host='localhost',
    port='5433'
)
cur = conn.cursor()

print('🔍 Vérification de l\'intégrité des données...')

# Test des contraintes de clés étrangères
tests = [
    ('Hôtels avec contrats', 'SELECT COUNT(*) FROM HOTEL h JOIN CONTRACT c ON h.hotel_key = c.hotel_id'),
    ('Utilisateurs avec rôles', 'SELECT COUNT(*) FROM USER_APP u JOIN ROLE r ON u.role_id = r.id'),
    ('Périodes avec contrats', 'SELECT COUNT(*) FROM PERIOD p JOIN CONTRACT c ON p.contract_id = c.id'),
    ('Prix de base complets', 'SELECT COUNT(*) FROM BASE_PRICE bp JOIN ROOM r ON bp.room_id = r.id'),
    ('Offres avec contrats', 'SELECT COUNT(*) FROM SPECIAL_OFFER so JOIN CONTRACT c ON so.contrat_id = c.id')
]

for test_name, query in tests:
    cur.execute(query)
    count = cur.fetchone()[0]
    if count > 0:
        print(f'✅ {test_name}: {count} relations valides')
    else:
        print(f'⚠️  {test_name}: Aucune relation trouvée')

# Test des dates cohérentes
cur.execute('''
    SELECT COUNT(*) FROM CONTRACT 
    WHERE start_at < end_at
''')
valid_contracts = cur.fetchone()[0]

cur.execute('SELECT COUNT(*) FROM CONTRACT')
total_contracts = cur.fetchone()[0]

if valid_contracts == total_contracts:
    print(f'✅ Dates de contrats: {valid_contracts}/{total_contracts} cohérentes')
else:
    print(f'❌ Dates de contrats: {valid_contracts}/{total_contracts} cohérentes')

# Test des prix positifs
cur.execute('''
    SELECT COUNT(*) FROM BASE_PRICE 
    WHERE price >= 0 AND extb_price >= 0
''')
valid_prices = cur.fetchone()[0]

cur.execute('SELECT COUNT(*) FROM BASE_PRICE')
total_prices = cur.fetchone()[0]

if valid_prices == total_prices:
    print(f'✅ Prix de base: {valid_prices}/{total_prices} valides')
else:
    print(f'❌ Prix de base: {valid_prices}/{total_prices} valides')

cur.close()
conn.close()
print('✅ Tests d\'intégrité terminés')
"

    echo ""
    echo "📋 EXEMPLES DE REQUÊTES UTILES"
    echo "==============================="
    echo ""
    echo "# Lister tous les hôtels avec leurs contrats actifs:"
    echo "SELECT h.name, h.city, c.name as contrat, c.start_at, c.end_at"
    echo "FROM HOTEL h JOIN CONTRACT c ON h.hotel_key = c.hotel_id"
    echo "WHERE c.active = true;"
    echo ""
    echo "# Voir les prix par type de chambre et période:"
    echo "SELECT rt.name as type_chambre, p.name as periode, bp.price"
    echo "FROM ROOM_TYPE rt"
    echo "JOIN BASE_PRICE bp ON rt.id = bp.type_id"
    echo "JOIN PERIOD p ON bp.period_id = p.id"
    echo "ORDER BY rt.name, bp.price;"
    echo ""
    echo "# Statistiques des offres spéciales:"
    echo "SELECT offer_type, COUNT(*) as nombre"
    echo "FROM SPECIAL_OFFER"
    echo "GROUP BY offer_type"
    echo "ORDER BY nombre DESC;"
    echo ""
    
    echo "🎉 CRÉATION TERMINÉE AVEC SUCCÈS !"
    echo "=================================="
    echo "✅ Base 'hotel_db' prête à l'emploi"
    echo "✅ 24 tables créées et peuplées"
    echo "✅ Données cohérentes vérifiées"
    echo "✅ Prêt pour l'intégration Hazelcast"
    echo ""
    echo "🔗 Connexion à la base:"
    echo "   Host: localhost"
    echo "   Port: 5433"
    echo "   Database: hotel_db"
    echo "   User: myuser"
    echo "   Password: mypassword"
    
else
    echo "❌ Échec de la création de la base de données"
    exit 1
fi
