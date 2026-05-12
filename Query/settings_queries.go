package Query 

const (
    UpsertBankingDetailsQuery = `
        INSERT INTO "BankingDetails" (
            "BankName", "AccountNumber", "ifscCode", "BankAddress", "LogoURL", 
            "AccountType", "SwiftCode", "userid", 
            "CreatedAt", "CreatedBy", "UpdatedAt", "UpdatedBy"
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), $9, NOW(), $9)
        ON CONFLICT (userid) 
        DO UPDATE SET 
            "BankName" = EXCLUDED."BankName",
            "AccountNumber"= EXCLUDED."AccountNumber",
            "ifscCode" = EXCLUDED."ifscCode",
            "BankAddress" = EXCLUDED."BankAddress",
            "LogoURL" = EXCLUDED."LogoURL",
            "AccountType" = EXCLUDED."AccountType",
            "SwiftCode" = EXCLUDED."SwiftCode",
            "UpdatedAt" = NOW(),
            "UpdatedBy" = EXCLUDED."UpdatedBy"
        RETURNING "DetailsID"`


    InsertBankingDetailsQuery = ` 
    INSERT INTO "BankingDetails" (
    "BankName", "AccountNumber", "ifscCode", "BankAddress", "LogoURL", 
    "AccountType", "SwiftCode", "userid", 
    "CreatedAt", "CreatedBy", "UpdatedAt", "UpdatedBy"
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), $9, NOW(), $9)
RETURNING "DetailsID";
    ` 

    UpdateBankingDetailsQuery = `
    UPDATE "BankingDetails"
    SET 
    "BankName" = $1,
    "AccountNumber" = $2,
    "ifscCode" = $3,
    "BankAddress" = $4,
    "LogoURL" = $5,
    "AccountType" = $6,
    "SwiftCode" = $7,
    "UpdatedAt" = NOW(),
    "UpdatedBy" = $8
WHERE "DetailsID" = $9
RETURNING "DetailsID";
    `

   
GetBankingDetailsQuery = `
    SELECT "DetailsID", "BankName", "AccountNumber", "ifscCode", "BankAddress", "LogoURL", "AccountType", "SwiftCode" 
    FROM "BankingDetails" 
    WHERE "userid" = $1 AND "DeletedAt" IS NULL
    ORDER BY "CreatedAt" DESC` 


DeleteBankingDetailsQuery = `
    UPDATE "BankingDetails" 
    SET "DeletedAt" = NOW(), "DeletedBy" = $1 
    WHERE "DetailsID" = $2 AND "userid" = $3`
) 

const( CreateCustomFieldQuery = `
    INSERT INTO "CustomFieldDefinitions" 
    ("FieldLabel", "FieldType", "IsRequired", "CreatedBy", "CreatedAt")
     VALUES ($1, $2, $3, $4, NOW())
     RETURNING "FieldID";`

 GetAllCustomFieldsQuery = `
    SELECT "FieldID", 
    "FieldLabel", 
    "FieldType", 
    "IsRequired", 
    "CreatedAt", 
    "CreatedBy"
FROM "CustomFieldDefinitions"
WHERE "DeletedAt" IS NULL
ORDER BY "CreatedAt" DESC;`

 DeleteCustomFieldQuery = `
    UPDATE "CustomFieldDefinitions"
SET 
    "DeletedAt" = NOW(),
    "DeletedBy" = $2
WHERE "FieldID" = $1;`
)