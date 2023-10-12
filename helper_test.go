package nano

import "testing"

func TestUnderlineToCamel(t *testing.T) {
	x := UnderlineToCamel("GUILD_CREATE")
	if x != "GuildCreate" {
		t.Fatal("expected GuildCreate but got", x)
	}
	x = UnderlineToCamel("GUILD_MEMBER_UPDATE")
	if x != "GuildMemberUpdate" {
		t.Fatal("expected GuildMemberUpdate but got", x)
	}
	x = UnderlineToCamel("OPEN_FORUM_THREAD_CREATE")
	if x != "OpenForumThreadCreate" {
		t.Fatal("expected OpenForumThreadCreate but got", x)
	}
	x = UnderlineToCamel("AUDIO_OR_LIVE_CHANNEL_MEMBER_ENTER")
	if x != "AudioOrLiveChannelMemberEnter" {
		t.Fatal("expected AudioOrLiveChannelMemberEnter but got", x)
	}
}
