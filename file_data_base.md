file_data_base.md


123monta456
postgres@go-175:~$ psql
psql (15.13 (Debian 15.13-0+deb12u1))
Type "help" for help.

postgres=# \password postgres
Enter new password for user "postgres": 123mont@456


engine = create_engine("postgresql://postgres:123mont%40456@localhost:5432/ma_base")

123Mont@456
information sur table film 
         column_name         |    data_type     
-----------------------------+------------------
 Rating                      | double precision
 Votes                       | bigint
 EpisodeDuration(in Minutes) | double precision
 Actors                      | text
 Title                       | text
 Years                       | text
 About                       | text
 Genres                      | text

SELECT column_name, data_type
FROM information_schema.columns
WHERE table_name = 'statistiques_api';

      column_name       |    data_type     
------------------------+------------------
 Communication_Skills   | bigint
 Projects_Completed     | bigint
 Prev_Sem_Result        | double precision
 CGPA                   | double precision
 Academic_Performance   | bigint
 IQ                     | bigint
 Extra_Curricular_Score | bigint
 Placement              | text
 Internship_Experience  | text
 College_ID             | text



pour voir la type de colonne 
\d+ film_rang

cette req pour voir le type 
SELECT 
    a.attname AS column_name,
    format_type(a.atttypid, a.atttypmod) AS data_type,
    a.attnotnull AS not_null,
    co.contype AS constraint_type
FROM 
    pg_attribute a
LEFT JOIN 
    pg_constraint co ON a.attrelid = co.conrelid AND a.attnum = ANY(co.conkey)
WHERE 
    a.attrelid = 'film_rang'::regclass
    AND a.attnum > 0 AND NOT a.attisdropped;


TRUNCATE TABLE statistiques_api RESTART IDENTITY;
