DO $$
DECLARE
    uid INT := 1;
    tid INT := 1;
BEGIN
    FOR tid IN 1..20 LOOP
        FOR i IN 1..10 LOOP
            INSERT INTO team_member (team_id, user_id)
            VALUES (tid, uid); 
            uid := uid + 1;
        END LOOP;
    END LOOP;
END $$;
