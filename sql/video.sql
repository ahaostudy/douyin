SELECT videos.*,
       CONCAT('http://192.168.1.9:8080/static/', videos.play_url)                   play_url,
       CONCAT('http://192.168.1.9:8080/static/', videos.cover_url)                  cover_url,
       (SELECT COUNT(*) FROM likes l WHERE videos.id = l.video_id)                  favorite_count,
       (SELECT COUNT(*) FROM comments c WHERE videos.id = c.video_id)               comment_count,
       EXISTS(SELECT * FROM likes l WHERE videos.id = l.video_id AND l.user_id = 1) is_favorite,
       users.*
FROM videos
         LEFT JOIN(SELECT u.*,
                          COUNT(DISTINCT v.id)  work_count,
                          COUNT(DISTINCT lv.id) total_favorited,
                          COUNT(DISTINCT lu.id) favorite_count
                   FROM users u
                            LEFT JOIN videos v ON u.id = v.author_id
                            LEFT JOIN likes lv ON v.id = lv.video_id
                            LEFT JOIN likes lu ON u.id = lu.user_id
                   GROUP BY u.id) users ON users.id = videos.author_id
WHERE videos.created_at <= '2023-08-10 18:13:33.56'
ORDER BY videos.created_at
LIMIT 30;

SELECT videos.*,
       CONCAT('http://192.168.1.9:8080/static/', videos.play_url)     play_url,
       CONCAT('http://192.168.1.9:8080/static/', videos.cover_url)    cover_url,
       (SELECT COUNT(*) FROM likes l WHERE videos.id = l.video_id)    favorite_count,
       (SELECT COUNT(*) FROM comments c WHERE videos.id = c.video_id) comment_count,
       true                                                           is_favorite,
       users.*
FROM videos
         LEFT JOIN(SELECT u.*,
                          COUNT(DISTINCT v.id)  work_count,
                          COUNT(DISTINCT lv.id) total_favorited,
                          COUNT(DISTINCT lu.id) favorite_count
                   FROM users u
                            LEFT JOIN videos v ON u.id = v.author_id
                            LEFT JOIN likes lv ON v.id = lv.video_id
                            LEFT JOIN likes lu ON u.id = lu.user_id
                   GROUP BY u.id) users ON users.id = videos.author_id
WHERE videos.created_at <= '2023-08-10 18:13:33.56'
  AND EXISTS(SELECT * FROM likes l WHERE videos.id = l.video_id AND l.user_id = 1)
ORDER BY videos.created_at
LIMIT 30;