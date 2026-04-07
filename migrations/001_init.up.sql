-- Создание таблицы пользователей
CREATE TABLE
    IF NOT EXISTS users (
        user_id SERIAL PRIMARY KEY,
        full_name VARCHAR(255) NOT NULL,
        middle_name VARCHAR(255),
        email VARCHAR(255) UNIQUE NOT NULL,
        password VARCHAR(255) NOT NULL
    );

-- Создание таблицы корзин
CREATE TABLE
    IF NOT EXISTS carts (
        cart_id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL REFERENCES users (user_id) ON DELETE CASCADE
    );

-- Создание таблицы товаров (продуктовый магазин)
CREATE TABLE
    IF NOT EXISTS products (
        product_id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        price DECIMAL(10, 2) NOT NULL,
        description TEXT,
        image VARCHAR(500),
        category VARCHAR(100),
        stock INTEGER DEFAULT 0
    );

-- Создание таблицы товаров в корзине
CREATE TABLE
    IF NOT EXISTS cart_items (
        cart_item_id SERIAL PRIMARY KEY,
        cart_id INTEGER NOT NULL REFERENCES carts (cart_id) ON DELETE CASCADE,
        product_id INTEGER NOT NULL REFERENCES products (product_id) ON DELETE CASCADE,
        quantity INTEGER NOT NULL DEFAULT 1,
        UNIQUE (cart_id, product_id)
    );

-- Индексы для производительности
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);

CREATE INDEX IF NOT EXISTS idx_carts_user_id ON carts (user_id);

CREATE INDEX IF NOT EXISTS idx_cart_items_cart_product ON cart_items (cart_id, product_id);

CREATE INDEX IF NOT EXISTS idx_products_category ON products (category);

CREATE INDEX IF NOT EXISTS idx_products_name ON products (name);

-- Добавляем тестового пользователя (Никита Крутой)
INSERT INTO
    users (full_name, middle_name, email, password)
VALUES
    (
        'Никита',
        'Крутой',
        'ivanov@example.ru',
        '7bT9xPqW'
    );

-- Добавляем тестового пользователя (Иван Иванов)
INSERT INTO
    users (full_name, middle_name, email, password)
VALUES
    ('Иван', 'Иванов', 'ivan@example.com', '123456');

-- Создаем корзины для существующих пользователей
INSERT INTO
    carts (user_id)
SELECT
    user_id
FROM
    users;

TRUNCATE TABLE products RESTART IDENTITY CASCADE;

-- Добавляем тестовые товары (продукты)
INSERT INTO
    products (name, price, description, image, category, stock)
VALUES
    -- Фрукты
    (
        'Яблоко Гренни Смит',
        150.00,
        'Сочные зеленые яблоки с кисло-сладким вкусом. Богаты витаминами и антиоксидантами.',
        'https://dixy.ru/upload/iblock/a3a/vicynai299r2wa3nszx5i5wct2hnjibc.webp',
        'Фрукты',
        50
    ),
    (
        'Банан',
        120.00,
        'Спелые сладкие бананы, источник калия и энергии.',
        'https://ir-3.ozone.ru/s3/multimedia-1-d/wc250/7851765829.jpg',
        'Фрукты',
        40
    ),
    (
        'Апельсин',
        180.00,
        'Сладкие сочные апельсины, богатые витамином C.',
        'https://ir-3.ozone.ru/s3/multimedia-1-g/wc1000/7557635356.jpg',
        'Фрукты',
        35
    ),
    (
        'Груша Конференция',
        210.00,
        'Сладкая и сочная груша с медовым ароматом.',
        'https://ir-3.ozone.ru/s3/multimedia-1-3/wc1000/8692265127.jpg',
        'Фрукты',
        30
    ),
    (
        'Киви',
        250.00,
        'Сочный киви, богатый витамином C и клетчаткой.',
        'https://ir-3.ozone.ru/s3/multimedia-1-t/wc500/7546973069.jpg',
        'Фрукты',
        45
    ),
    -- Овощи
    (
        'Помидоры черри',
        280.00,
        'Сладкие маленькие помидоры, идеальны для салатов.',
        'https://ir-3.ozone.ru/s3/multimedia-1-2/wc1000/9139825946.jpg',
        'Овощи',
        45
    ),
    (
        'Огурцы гладкие',
        130.00,
        'Хрустящие свежие огурцы, выращенные в теплицах.',
        'https://ir-3.ozone.ru/s3/multimedia-1-1/wc1000/9437859541.jpg',
        'Овощи',
        40
    ),
    (
        'Картофель молодой',
        55.00,
        'Молодой картофель, отлично подходит для варки и жарки.',
        'https://ir-3.ozone.ru/s3/multimedia-1-1/wc1000/7995100501.jpg',
        'Овощи',
        60
    ),
    (
        'Морковь',
        70.00,
        'Сладкая и хрустящая морковь, богатая бета-каротином.',
        'https://ir-3.ozone.ru/s3/multimedia-1-2/wc1000/7557622238.jpg',
        'Овощи',
        55
    ),
    (
        'Болгарский перец',
        190.00,
        'Сочный сладкий перец красного цвета.',
        'https://ir-3.ozone.ru/s3/multimedia-1-c/wc1000/8145573132.jpg',
        'Овощи',
        35
    ),
    -- Ягоды
    (
        'Клубника',
        350.00,
        'Сладкая и ароматная клубника, урожай сезона.',
        'https://ir-3.ozone.ru/s3/multimedia-1-a/wc1000/9444252718.jpg',
        'Ягоды',
        25
    ),
    (
        'Черника',
        400.00,
        'Свежая черника, богатая антиоксидантами.',
        'https://cdn-irec.r-99.com/sites/default/files/imagecache/300o/product-images/10297/fTfSiRxMumA2Udxb2yJegQ.jpeg',
        'Ягоды',
        20
    ),
    (
        'Малина',
        380.00,
        'Нежная и ароматная малина.',
        'https://cdn-irec.r-99.com/sites/default/files/imagecache/150o/product-images/10297/ld8VOyi4tUUnOGFeiHvMTQ.jpeg',
        'Ягоды',
        20
    ),
    (
        'Ежевика',
        420.00,
        'Сочная и полезная ежевика.',
        'https://cdn-irec.r-99.com/sites/default/files/imagecache/150o/product-images/599862/zwUcFQ0yXOt8H0ymGYSew.jpg',
        'Ягоды',
        18
    ),
    (
        'Красная смородина',
        220.00,
        'Кисло-сладкие ягоды красной смородины.',
        'https://cdn-irec.r-99.com/sites/default/files/imagecache/150o/product-images/10297/x6exoi4vXUgF7k6gdSnw.jpeg',
        'Ягоды',
        30
    );

-- Добавляем примеры товаров в корзины (опционально)
INSERT INTO
    cart_items (cart_id, product_id, quantity)
VALUES
    (1, 1, 2), -- Корзина Никиты: 2 кг яблок
    (1, 3, 3), -- Корзина Никиты: 3 апельсина
    (2, 5, 1), -- Корзина Ивана: 1 литр молока
    (2, 9, 1);

-- Корзина Ивана: 1 сыр пармезан