-- SQLite


INSERT INTO voiture(
    titre,
    marque,
    quantité,
    created_at
) VALUES (
    'tonm',
    'tongrap',
    4000,
    DATETIME('now')
);

INSERT INTO voiture(
    titre,
    marque,
    quantité,
    created_at
) VALUES (
    'tonp',
    'tongrap',
    5000,
    DATETIME('now')
);


SELECT *
FROM voiture
where quantité > 4500