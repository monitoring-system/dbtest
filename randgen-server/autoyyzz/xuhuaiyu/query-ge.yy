query:
    select;

select:
    SELECT a
    FROM (SELECT _field AS a, _field AS b
          FROM _table
    ) tmp
    WHERE filter
    ORDER BY a;

filter:
    a >= b
|   a >= -1
|   a >= 1
|   a >= 0
|   a >= null
|   a >= 4294967295
|   a >=  -2147483648
|   a >=  2147483647
|   a >=  -9223372036854775808
|   a >=  9223372036854775807
|   a >=  18446744073709551615