-- Создать таблицу groups, если она ещё не существует
CREATE TABLE IF NOT EXISTS groups (
                                      id SERIAL PRIMARY KEY,
                                      "group" VARCHAR(255) NOT NULL
    );

-- Создать таблицу songs
CREATE TABLE songs (
                       id SERIAL PRIMARY KEY,
                       group_id INT REFERENCES groups(id) ON DELETE CASCADE NOT NULL,
                       song VARCHAR(255) NOT NULL,
                       text TEXT NOT NULL,
                       release_date VARCHAR(255) NOT NULL,
                       link VARCHAR(255) NOT NULL
);

-- Создать индексы для songs
CREATE INDEX IF NOT EXISTS idx_songs_group_id ON songs (group_id);
CREATE INDEX IF NOT EXISTS idx_songs_song ON songs (song);
CREATE INDEX IF NOT EXISTS idx_songs_release ON songs (release_date);
CREATE INDEX IF NOT EXISTS idx_songs_link ON songs (link);

-- Создать индексы для groups
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'ind_groups_name') THEN
CREATE INDEX ind_groups_name ON groups ("group");
END IF;
END $$;
