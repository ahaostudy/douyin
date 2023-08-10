# 查询用户信息及作品数、获赞数、喜欢数
SELECT
    u.*,
    COUNT(DISTINCT v.id) work_count,
    COUNT(DISTINCT lv.id) total_favorited,
    COUNT(DISTINCT lu.id) favorite_count
FROM users u
LEFT JOIN videos v ON u.id = v.author_id
LEFT JOIN likes lv ON v.id = lv.video_id
LEFT JOIN likes lu on u.id = lu.user_id
WHERE u.id = 15;
