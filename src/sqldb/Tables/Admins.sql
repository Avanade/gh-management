CREATE TABLE [dbo].[Admins]
(
  [UserPrincipalName] VARCHAR(100) NOT NULL PRIMARY KEY
  CONSTRAINT [FK_Admins_Users] FOREIGN KEY (UserPrincipalName) REFERENCES Users(UserPrincipalName), 
)
