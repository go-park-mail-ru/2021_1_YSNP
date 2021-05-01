#!/usr/bin/env tarantool
json=require('json')

-- Настроить базу данных
box.cfg {
    listen = 3301
}

-- При поднятии БД создаем спейсы и индексы
box.once('init', function()
    box.schema.space.create('sessions')
    box.schema.space.create('trends')
    box.schema.space.create('trends_products')
    box.schema.user.passwd('pass')
    box.space.sessions:create_index('primary',
        { type = 'TREE', parts = {1, 'string'}})
    box.space.trends:create_index('primary',
        { type = 'TREE', parts = {1, 'number'}})
    box.space.trends_products:create_index('primary',
        { type = 'TREE', parts = {1, 'number'}})

end)

function check_session(session_id)
    local session_id = box.space.sessions:select{session_id}[1]
    print('found session', session_id)
    return session_id
end

function get_user_trend(userID)
    local value = box.space.trends:select{userID}[1][2]
    return value
end

function get_user_trends_products(userID)
    local value = box.space.trends_products:select{userID}[1][2]
    return value
end
