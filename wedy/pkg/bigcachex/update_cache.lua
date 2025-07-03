local key = KEYS[1]
local version = tonumber(redis.call('GET', key))
local data=ARGV[2]

if version >= tonumber(ARGV[1]) then
    -- update data
    redis.call('SET', key, data)
    return 1 -- 操作成功
end

return -1  -- skip update
