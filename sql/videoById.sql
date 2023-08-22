SELECT videos.*,
       CONCAT('http://172.30.24.3:8080/static/',videos.play_url)                     play_url,
       CONCAT('http://172.30.24.3:8080/static/',videos.cover_url)                    cover_url,
       (SELECT COUNT(*) FROM likes l WHERE videos.id=l.video_id)            favorite_count,
       (SELECT COUNT(*) FROM comments c WHERE videos.id=c.video_id)         comment_count,
       EXISTS(SELECT * FROM likes l WHERE videos.id = l.video_id AND l.user_id = 1)  is_favorite,
       users.*
FROM videos
         LEFT JOIN (
    SELECT u.*,
           COUNT(DISTINCT v.id) work_count,
           COUNT(DISTINCT vl.id) total_favorited,
           COUNT(DISTINCT ul.id) favorite_count
    FROM users u
             LEFT JOIN videos v ON u.id = v.authorID
             LEFT JOIN likes vl ON v.id = vl.video_id
             LEFT JOIN likes ul ON u.id = ul.user_id
    GROUP BY u.id
)  users ON users.id = videos.authorID
WHERE videos.authorID = 1
ORDER BY videos.created_at;