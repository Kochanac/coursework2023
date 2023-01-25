create table test_1 (
    id int8,
    data int8
);

create index test_1_idx on test_1 using hash(id);

insert into test_1(id, data)
select i + 1000000000 * (i%8),
       i + 4824246592 + 1000000000 * (i%8)
from generate_series(6978947857221617788::int8, 6978947857221617788::int8 + 1e6) AS t(i);

