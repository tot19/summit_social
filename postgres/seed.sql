-- seed.sql

-- Insert sample posts
INSERT INTO posts (id, content, image, likes, poster, comments_count)
VALUES 
(1, 'Welcome to Summit Social!', '/images/welcome.png', 0, 'admin', 0),
(2, 'Check out this awesome picture!', '/images/pic1.png', 10, 'user1', 2),
(3, 'Another amazing post!', '/images/post2.png', 5, 'user2', 3)
ON CONFLICT (id) DO NOTHING;

