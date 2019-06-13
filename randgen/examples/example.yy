query:
    select ;

select:
    SELECT coalesce
    FROM _table
    WHERE condition ;

coalesce:
    COALESCE( _field , 0) | COALESCE( _field_list ) ;

condition:
    _field IS NULL | _field = 1111 | _field = 'hello' ;