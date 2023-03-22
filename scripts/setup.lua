-- initialization script for REAL tarantool server only
box.cfg{listen = 3301}
s = box.schema.space.create('tester')
s:format({
         {name = 'id', type = 'unsigned'},
         {name = 'band_name', type = 'string'},
         {name = 'year', type = 'unsigned'}
         })
s:create_index('primary', {
         type = 'hash',
         parts = {'id'}
         })
s:create_index('secondary', {
         type = 'hash',
         parts = {'band_name'}
         })
s:insert{1, 'Roxette', 1986}
s:insert{2, 'Scorpions', 2015}
s:insert{3, 'Ace of Base', 1993}
