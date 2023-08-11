# 查询用户信息及作品数、获赞数、喜欢数
SELECT u.*,
       COUNT(DISTINCT v.id)                                                         work_count,
       COUNT(DISTINCT l.id)                                                         total_favorited,
       (SELECT COUNT(*) FROM likes l WHERE l.user_id = u.id)                        favorite_count,
       (SELECT COUNT(*) FROM follows f WHERE f.follower_id = u.id)                  follow_count,
       (SELECT COUNT(*) FROM follows f WHERE f.user_id = u.id)                      follower_count,
       EXISTS(SELECT * FROM follows f WHERE f.follower_id = u.id AND f.user_id = 5) is_follow
FROM users u
         LEFT JOIN videos v ON u.id = v.author_id
         LEFT JOIN likes l ON v.id = l.video_id
group by u.id;

SELECT u.*,
       COUNT(DISTINCT v.id)                                                         work_count,
       COUNT(DISTINCT l.id)                                                         total_favorited,
       (SELECT COUNT(*) FROM likes l WHERE l.user_id = u.id)                        favorite_count,
       (SELECT COUNT(*) FROM follows f WHERE f.follower_id = u.id)                  follow_count,
       (SELECT COUNT(*) FROM follows f WHERE f.user_id = u.id)                      follower_count,
       EXISTS(SELECT * FROM follows f WHERE f.follower_id = u.id AND f.user_id = 1) is_follow
FROM users u
         LEFT JOIN videos v ON u.id = v.author_id
         LEFT JOIN likes l ON v.id = l.video_id
WHERE u.id IN (SELECT user_id FROM follows WHERE follower_id = 1)
GROUP BY u.id;