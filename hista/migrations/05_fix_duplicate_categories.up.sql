with duplicates as (
	select 
        name, 
        min(id) as targetId 
    from symptom_categories 
    group by name 
    having count(*) > 1
)

update symptoms 
set symptom_category_id = targets.targetId
from (
	select * from symptom_categories c join duplicates d on c.name = d.name
) as targets
where symptoms.symptom_category_id = targets.id;

 
delete from symptom_categories where id in (
	select c.id 
	from symptoms s right join symptom_categories c on s.symptom_category_id = c.id
	where s.id is null
);