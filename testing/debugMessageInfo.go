/*fmt.Println("GetConversation : ", v.Message.GetConversation())
fmt.Println("Sender : ", v.Info.Sender)
fmt.Println("Sender Number : ", v.Info.Sender.User)
fmt.Println("IsGroup : ", v.Info.IsGroup)
fmt.Println("MessageSource : ", v.Info.MessageSource)
fmt.Println("ID : ", v.Info.ID)
fmt.Println("PushName : ", v.Info.PushName)
fmt.Println("BroadcastListOwner : ", v.Info.BroadcastListOwner)
fmt.Println("Category : ", v.Info.Category)
fmt.Println("Chat : ", v.Info.Chat)
fmt.Println("DeviceSentMeta : ", v.Info.DeviceSentMeta)
fmt.Println("IsFromMe : ", v.Info.IsFromMe)
fmt.Println("MediaType : ", v.Info.MediaType)
fmt.Println("Multicast : ", v.Info.Multicast)
fmt.Println("Info.Chat.Server : ", v.Info.Chat.Server)
if v.Info.Chat.Server == "g.us" {
	groupInfo, err := client.GetGroupInfo(v.Info.Chat)
	fmt.Println("error GetGroupInfo : ", err)
	fmt.Println("Nama Group : ", groupInfo.GroupName.Name)
}*/