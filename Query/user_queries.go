package Query

const (
	GetAllUsersQuery = `
		SELECT 
			"userid", 
			"usercode", 
			"username", 
			COALESCE("firstname", '') AS "firstname", 
			COALESCE("lastname", '') AS "lastname", 
			"roleid" 
		FROM active_users;`

	CreateUserQuery = `
    INSERT INTO users (
        usercode,
        username,
        password,
        password_plain,
        firstname,
        lastname,
        roleid,
        emailid,
        email_plain,
        is_first_login
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    RETURNING userid;`

	GetUserByEmailQuery = `
    SELECT 
        u.userid, u.password, u.isactive, r.rolename, u.username, u.is_first_login 
    FROM users u
    JOIN roles r ON u.roleid = r.roleid
    WHERE u.email_plain = $1;`

	UpdateUserQuery = `
    UPDATE "users" 
    SET 
        "username" = $1, 
        "firstname" = $2, 
        "lastname" = $3, 
        "roleid" = $4,
        "emailid" = $5,
        "updatedat" = NOW()
    WHERE "userid" = $6;
`

	DeleteUserQuery = `
    UPDATE "users" 
    SET 
        "isactive" = false, 
        "deletedat" = NOW(),
        "deletedby" = $1
    WHERE "userid" = $2;`

	GetUserByIDQuery = `
    SELECT userid, usercode, username, firstname, lastname, roleid, emailid
    FROM active_users 
    WHERE userid = $1 LIMIT 1;`

	GetUserProfileByIDQuery = `
    SELECT u.firstname, u.lastname, r.rolename 
    FROM users u
    JOIN roles r ON u.roleid = r.roleid
    WHERE u.userid = $1;`

	ResetPasswordQuery = `
UPDATE users
SET
    password = $1,
    password_plain = $2,
    is_first_login = false
WHERE userid = $3;`
    


)