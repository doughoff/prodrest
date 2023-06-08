begin;

alter table "entities"
    add constraint ci_unique unique(ci),
    add constraint ruc_unique unique(ruc)
;

alter table "stock_movements"
alter column entity_id drop not null;

commit;