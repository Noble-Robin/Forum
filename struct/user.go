package struct

type user struct{
	user_id string
	username string
	name string
	email string
	password string
}

type user_conversation struct {
	conversation_id string
	user_id users
	conversation_id conversations

}

type message struct{
	id_message string
	id_conv conversations
	id_user users
	contenu string
	creation Timestamp
}

type groups struct{
	id string
	groupname string
	creation Timestamp

}

type groups_members struct{
	group_id groups 
	user_id users
}

type conversation struct{
	id string
	creation TIMESTAMP
}

type conversation_participants struct{
	conversation_id conversations
	user_id users
}

type contacts struct{
	id string
	user_id users
	status string
}

