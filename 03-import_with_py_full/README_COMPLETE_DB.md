# 🏨 Base de Données Hôtelière Complète - Documentation

## 📋 Vue d'ensemble

Ce projet crée et peuple une base de données PostgreSQL complète pour un système de gestion hôtelière avec **24 tables interconnectées**. La base est conçue pour gérer tous les aspects d'un système de réservation hôtelière professionnel.

## 🗂️ Structure des fichiers

```
hotel_hazelcast/import_with_py/
├── base.py                     # 📄 Version originale (6 tables)
├── complete_hotel_db.py        # 📄 Version complète (24 tables) ⭐
├── test_complete_db.sh         # 🧪 Script de test et vérification
└── README_COMPLETE_DB.md       # 📖 Cette documentation
```

## 🏗️ Architecture de la base de données

### 📊 Répartition des 24 tables

| Catégorie | Tables | Description |
|-----------|--------|-------------|
| **Base** | 6 tables | COUNTRY, CITY, HOTEL, USER_APP, ROLE, CONTRACT |
| **Structure** | 4 tables | ROOM, ROOM_TYPE, BOARD, PERIOD |
| **Liaison** | 2 tables | CONTRACT_ROOM, PERIOD_BOARD |
| **Tarification** | 4 tables | BASE_PRICE, BOARD_PRICES, ROOM_PRICES, ROOM_TYPE_PRICES |
| **Promotion** | 4 tables | SPECIAL_OFFER, DAY_PROMOTION, ROOM_PROMOTION, PAX_PROMOTION |
| **Gestion** | 4 tables | SUPPLEMENT, CHILD_PRICE, CANCELLATION_CONDITION, ALLOTMENT |

### 🔗 Schéma de dépendances

```
COUNTRY ──┐
          │
CITY ─────┼──► HOTEL ──► CONTRACT ──┬──► PERIOD ──┬──► BASE_PRICE
          │                        │             │
ROLE ─────┴──► USER_APP ────────────┘             ├──► BOARD_PRICES
                                                  │
ROOM_TYPE ────────────────────────────────────────┼──► ROOM_TYPE_PRICES
                                                  │
ROOM ─────────────────────────────────────────────┼──► ROOM_PRICES
                                                  │
BOARD ────────────────────────────────────────────┴──► PERIOD_BOARD
                                                  
CONTRACT ──► SPECIAL_OFFER ──┬──► DAY_PROMOTION
                             ├──► ROOM_PROMOTION
                             └──► PAX_PROMOTION

PERIOD ──┬──► SUPPLEMENT
         ├──► CHILD_PRICE
         ├──► CANCELLATION_CONDITION
         └──► ALLOTMENT
```

## 🚀 Utilisation

### 📋 Prérequis

```bash
# Dépendances Python
pip install psycopg2-binary faker

# PostgreSQL accessible sur :
Host: localhost
Port: 5433
User: myuser
Password: mypassword
```

### ▶️ Exécution

#### Méthode 1 : Script automatique (recommandé)
```bash
cd /workspaces/mlops-zoom/hotel_hazelcast/import_with_py/
./test_complete_db.sh
```

#### Méthode 2 : Exécution manuelle
```bash
python3 complete_hotel_db.py
```

### 📊 Résultat attendu

```
🎯 TOTAL                     :   XXX enregistrements
🏗️ TABLES CRÉÉES             :    24 tables

✅ Base de données hôtelière créée avec succès
✅ 24 tables créées et peuplées
✅ Données cohérentes avec contraintes respectées
✅ Prêt pour l'intégration avec Hazelcast
```

## 📋 Détail des tables

### 🏛️ Tables de base (6)

#### 1. **COUNTRY** - Pays
```sql
- id (PK)          : Identifiant unique
- name             : Nom du pays
- code             : Code ISO (FR, ES, IT...)
```

#### 2. **CITY** - Villes
```sql
- id (PK)          : Identifiant unique
- citycode         : Code de la ville
- name             : Nom de la ville
- countryname      : Nom du pays
- countryid (FK)   : Référence au pays
```

#### 3. **HOTEL** - Hôtels
```sql
- hotel_key (PK)   : Clé unique (HOT001, HOT002...)
- name             : Nom de l'hôtel
- city             : Ville de l'hôtel
- country          : Pays de l'hôtel
- stars            : Nombre d'étoiles (1-5)
- address          : Adresse complète
- mail             : Email de contact
- latitude/longitude : Coordonnées GPS
- phone            : Téléphone
- description      : Description longue
- short_description: Description courte
```

#### 4. **ROLE** - Rôles utilisateurs
```sql
- id (PK)          : Identifiant unique
- name             : Nom du rôle (Admin, Agent, Client...)
```

#### 5. **USER_APP** - Utilisateurs
```sql
- id (PK)          : Identifiant unique
- status           : Statut (actif, suspendu...)
- marge            : Marge générale
- marge_operation  : Marge opérationnelle
- solde            : Solde disponible
- solde_rouge      : Découvert autorisé
- currency         : Devise (EUR, USD...)
- maxrequest       : Limite de requêtes
- group            : Groupe d'appartenance
- marge_b2b        : Marge B2B
- marge_xml        : Marge XML/API
- role_id (FK)     : Référence au rôle
```

#### 6. **CONTRACT** - Contrats
```sql
- id (PK)          : Identifiant unique
- name             : Nom du contrat
- hotel_id (FK)    : Référence à l'hôtel
- start_at         : Date de début
- end_at           : Date de fin
- access           : Type d'accès
- active           : Contrat actif ?
- currency         : Devise
- market           : Marché ciblé
- client_id (FK)   : Référence au client
```

### 🏗️ Tables de structure (4)

#### 7. **ROOM_TYPE** - Types de chambres
```sql
- id (PK)          : Identifiant unique
- name             : Nom du type (Single, Double...)
- code             : Code abrégé (SS, DS...)
- client_id (FK)   : Client spécifique
```

#### 8. **ROOM** - Chambres
```sql
- id (PK)          : Identifiant unique
- name             : Nom de la chambre
- name_code        : Code de la chambre
- client_id (FK)   : Client associé
- max_pax          : Maximum de personnes
- min_pax          : Minimum de personnes
- child            : Enfants permis
- min_adult        : Minimum d'adultes
- max_adult        : Maximum d'adultes
```

#### 9. **BOARD** - Types de pension
```sql
- id (PK)          : Identifiant unique
- name             : Nom (Room Only, B&B, Half Board...)
- definition       : Description détaillée
- client_id (FK)   : Client spécifique
```

#### 10. **PERIOD** - Périodes de contrat
```sql
- id (PK)          : Identifiant unique
- contract_id (FK) : Référence au contrat
- start_at         : Date de début
- end_at           : Date de fin
- minimum_stay     : Séjour minimum
- delai            : Délai de réservation
- min_stay         : Séjour minimum autorisé
- name             : Nom de la période
- code             : Code de la période
```

### 🔗 Tables de liaison (2)

#### 11. **CONTRACT_ROOM** - Liaison contrat/chambre
```sql
- id (PK)                    : Identifiant unique
- contract_id (FK)           : Référence au contrat
- room_id (FK)               : Référence à la chambre
- room_type_id (FK)          : Référence au type
- maximum_number_people      : Max personnes
- minimum_number_people      : Min personnes
- max_adult/min_adult        : Limites adultes
- child                      : Enfants
- code_room/code_type        : Codes associés
```

#### 12. **PERIOD_BOARD** - Liaison période/pension
```sql
- id (PK)            : Identifiant unique
- period_id (FK)     : Référence à la période
- board_id (FK)      : Référence à la pension
- extra_board_id (FK): Pension supplémentaire
```

### 💰 Tables de tarification (4)

#### 13. **BASE_PRICE** - Prix de base
```sql
- id (PK)          : Identifiant unique
- room_id (FK)     : Référence à la chambre
- period_id (FK)   : Référence à la période
- type_id (FK)     : Référence au type
- board_id (FK)    : Référence à la pension
- price            : Prix de base
- par_pax          : Prix par personne ?
- extb_price       : Prix lit supplémentaire
- operation        : Type d'opération
```

#### 14. **BOARD_PRICES** - Prix des pensions
#### 15. **ROOM_PRICES** - Prix par chambre
#### 16. **ROOM_TYPE_PRICES** - Prix par type de chambre

### 🎁 Tables de promotion (4)

#### 17. **SPECIAL_OFFER** - Offres spéciales
```sql
- id (PK)                      : Identifiant unique
- receive_date                 : Date de réception
- requestdate_from/to          : Période de requêtes
- checkin_date/checkout_date   : Période de séjour
- min_stay/max_stay            : Durée de séjour
- order                        : Priorité d'affichage
- offer_type                   : Type d'offre
- origine                      : Source de l'offre
- definition                   : Description
- note                         : Remarques
- bookingdate_range_exclusif   : Plage exclusive ?
- in_contrat                   : Associée au contrat ?
- priority                     : Offre prioritaire ?
- contrat_id (FK)              : Référence au contrat
- client_id (FK)               : Référence au client
```

#### 18. **DAY_PROMOTION** - Promotions par jour
#### 19. **ROOM_PROMOTION** - Promotions par chambre
#### 20. **PAX_PROMOTION** - Promotions par nombre de personnes

### ⚙️ Tables de gestion (4)

#### 21. **SUPPLEMENT** - Suppléments tarifaires
#### 22. **CHILD_PRICE** - Prix enfants
#### 23. **CANCELLATION_CONDITION** - Conditions d'annulation
#### 24. **ALLOTMENT** - Quotas de chambres

## 📊 Données générées

### 📈 Quantités par défaut

| Table | Enregistrements | Description |
|-------|----------------|-------------|
| COUNTRY | 10 | Pays européens + autres |
| CITY | 20 | Villes touristiques |
| HOTEL | 12 | Hôtels 1-5 étoiles |
| USER_APP | 15 | Utilisateurs divers |
| ROLE | 5 | Rôles standards |
| CONTRACT | 15 | Contrats actifs/inactifs |
| ROOM_TYPE | 8 | Types de chambres |
| ROOM | 20 | Chambres individuelles |
| BOARD | 6 | Pensions RO à AI Premium |
| PERIOD | 25 | Périodes variées |
| ... | ... | Autres tables peuplées |

### 🎲 Caractéristiques des données

- **Données réalistes** : Faker pour noms, adresses, emails
- **Cohérence** : Contraintes de clés étrangères respectées
- **Variété** : Différents types, statuts, devises
- **Intégrité** : Dates cohérentes, prix positifs
- **Relations** : Liens logiques entre les entités

## 🔬 Tests et vérifications

### ✅ Tests automatiques inclus

1. **Test de connexion** PostgreSQL
2. **Création des tables** dans l'ordre correct
3. **Insertion des données** avec gestion d'erreurs
4. **Vérification d'intégrité** des contraintes
5. **Tests de cohérence** des dates et prix
6. **Statistiques complètes** par table

### 🧪 Exemples de requêtes

```sql
-- Hôtels avec contrats actifs
SELECT h.name, h.city, c.name as contrat, c.start_at, c.end_at
FROM HOTEL h JOIN CONTRACT c ON h.hotel_key = c.hotel_id
WHERE c.active = true;

-- Prix par type de chambre et période
SELECT rt.name as type_chambre, p.name as periode, bp.price
FROM ROOM_TYPE rt
JOIN BASE_PRICE bp ON rt.id = bp.type_id
JOIN PERIOD p ON bp.period_id = p.id
ORDER BY rt.name, bp.price;

-- Statistiques des offres spéciales
SELECT offer_type, COUNT(*) as nombre
FROM SPECIAL_OFFER
GROUP BY offer_type
ORDER BY nombre DESC;
```

## 🔄 Intégration avec Hazelcast

Cette base est **parfaitement compatible** avec votre API Hazelcast existante :

1. **Structures identiques** : Les tables principales correspondent aux modèles Go
2. **Clés cohérentes** : hotel_key, ID appropriés
3. **Données réalistes** : Prêtes pour les tests et la production
4. **Extensibilité** : Facile d'ajouter de nouvelles tables

### 🔗 Prochaines étapes

1. **Export vers Hazelcast** : Script pour charger les données
2. **Synchronisation** : Mise à jour en temps réel
3. **API REST** : Endpoints pour toutes les tables
4. **Interface admin** : Gestion complète des données

## 📝 Résumé

✅ **24 tables créées** avec relations complètes  
✅ **Données cohérentes** et réalistes générées  
✅ **Tests automatiques** inclus  
✅ **Documentation complète** fournie  
✅ **Compatible Hazelcast** pour l'intégration  
✅ **Prêt pour la production** hôtelière  

🎯 **Base de données professionnelle complète pour système hôtelier !**
