resource onelogin_roles role1 {
    name = "role_1_acctest"
}

resource onelogin_roles role2 {
    name = "role_2_acctest"
}

resource onelogin_users user1 {
    username = "testy.mctesterson.acctest"
    email = "testy.mctesterson.acctest@example.com"
}

resource onelogin_users user2 {
    username = "boaty.mcboatface.acctest"
    email = "boaty.mcboatface.acctest@example.com"
}

resource onelogin_privileges super_admin {
  name = "super admin"
  description = "description"
  user_ids = [tonumber(onelogin_users.user1.id), tonumber(onelogin_users.user2.id)]
  role_ids = [tonumber(onelogin_roles.role1.id)]
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