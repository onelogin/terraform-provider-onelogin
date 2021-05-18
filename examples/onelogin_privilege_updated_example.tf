resource onelogin_roles role1 {
    name = "role_1"
}

resource onelogin_roles role2 {
    name = "role_1"
}

resource onelogin_users user1 {
    username = "testy.mctesterson"
    email = "testy.mctesterson@onelogin.com"
}

resource onelogin_users user2 {
    username = "boaty.mcboatface"
    email = "boaty.mcboatface@onelogin.com"
}

resource onelogin_privileges super_admin {
  name = "super admin"
  description = "description"
  user_ids = [user1.id]
  role_ids = [role2.id]
  privilege {
	statement {
		effect = "Allow"
		action = ["apps:List"]
		scope = ["*"]
	}
	statement {
		effect = "Allow"
		action = ["users:List"]
		scope = ["*"]
	}
  }
}