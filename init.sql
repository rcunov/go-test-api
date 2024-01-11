DROP TABLE IF EXISTS albums;

CREATE TABLE albums (
    id INTEGER PRIMARY KEY AUTOINCREMENT, 
    title VARCHAR NOT NULL, 
    artist VARCHAR NOT NULL, 
    price DECIMAL(5,2) NOT NULL
);

INSERT INTO albums
    (title, artist, price)
VALUES
    ('Blue Train', 'John Coltrane', 56.99),
    ('Bleed the Future', 'AUM', 19.99),
    ('Super Hexagon', 'Chipzel', 8.0),
    ('Hirschbrunnen', 'delving', 14.99);

-- .mode column
-- .headers on
