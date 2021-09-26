select s.description
     , c.description
     , m.ts
     , m.total_probable_cases
     , m.active_cases
     , m.total_positive_cases
     , m.resolved
from hasd_covid.metric m
   , hasd_covid.school s
   , hasd_covid.category c
where 1 = 1
  and m.school_skey = s.school_skey
  and m.category_skey = c.category_skey
  #   and ts > (sysdate() - interval '7' DAY)
order by ts, s.description, c.description desc;


