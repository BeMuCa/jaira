package tui

// Code generated from icon/jAIra.png by scripts/iconpreview. Do not edit.
//
// The icon is pre-rendered rather than decoded at startup. The board is glanced
// at constantly and instant startup is a stated constraint, so paying for a PNG
// decode and a flood fill on every launch to produce a fixed picture would be
// the wrong trade.
//
// This is the crescent alone, with the prompt glyphs and rays removed and the
// background keyed out so the terminal shows through. The colours are the
// image's own, which are near-white: on a light-background terminal the moon
// will have poor contrast, and the silhouette style is the fix if that matters.
const iconArt = "           \x1b[38;2;82;80;93m\x1b[48;2;236;232;254m▀\x1b[0m\x1b[38;2;226;223;241m▀\x1b[0m                       \n" +
	"         \x1b[38;2;99;97;109m\x1b[48;2;235;232;252m▀\x1b[0m\x1b[38;2;229;224;255m\x1b[48;2;240;236;250m▀\x1b[0m\x1b[38;2;86;84;96m▀\x1b[0m                        \n" +
	"       \x1b[38;2;246;243;255m▄\x1b[0m\x1b[38;2;228;224;250m\x1b[48;2;239;236;255m▀\x1b[0m\x1b[38;2;235;231;254m\x1b[48;2;239;236;255m▀\x1b[0m\x1b[38;2;96;94;109m▀\x1b[0m                         \n" +
	"      \x1b[38;2;91;89;103m▄\x1b[0m\x1b[38;2;238;235;254m\x1b[48;2;238;235;255m▀\x1b[0m\x1b[38;2;239;236;254m\x1b[48;2;241;239;252m▀\x1b[0m\x1b[38;2;241;239;251m\x1b[48;2;244;242;255m▀\x1b[0m                          \n" +
	"      \x1b[38;2;231;229;243m\x1b[48;2;233;230;247m▀\x1b[0m\x1b[38;2;239;236;254m\x1b[48;2;240;237;253m▀\x1b[0m\x1b[38;2;242;240;254m\x1b[48;2;242;241;254m▀\x1b[0m\x1b[38;2;249;247;255m\x1b[48;2;246;244;255m▀\x1b[0m                          \n" +
	"      \x1b[38;2;238;235;253m\x1b[48;2;239;236;253m▀\x1b[0m\x1b[38;2;240;237;253m\x1b[48;2;241;239;254m▀\x1b[0m\x1b[38;2;243;243;254m\x1b[48;2;243;243;253m▀\x1b[0m\x1b[38;2;246;245;254m\x1b[48;2;245;243;254m▀\x1b[0m                          \n" +
	"      \x1b[38;2;238;235;252m\x1b[48;2;94;92;105m▀\x1b[0m\x1b[38;2;241;239;254m\x1b[48;2;241;239;254m▀\x1b[0m\x1b[38;2;244;244;254m\x1b[48;2;243;242;255m▀\x1b[0m\x1b[38;2;246;245;253m\x1b[48;2;245;243;253m▀\x1b[0m\x1b[38;2;115;112;123m\x1b[48;2;245;244;252m▀\x1b[0m                         \n" +
	"       \x1b[38;2;238;235;251m\x1b[48;2;240;237;254m▀\x1b[0m\x1b[38;2;242;241;254m\x1b[48;2;242;240;254m▀\x1b[0m\x1b[38;2;245;243;253m\x1b[48;2;244;244;254m▀\x1b[0m\x1b[38;2;248;248;255m\x1b[48;2;247;246;253m▀\x1b[0m\x1b[38;2;188;186;196m\x1b[48;2;249;248;253m▀\x1b[0m\x1b[38;2;133;128;147m▄\x1b[0m                       \n" +
	"       \x1b[38;2;100;98;109m▀\x1b[0m\x1b[38;2;239;237;253m\x1b[48;2;232;229;245m▀\x1b[0m\x1b[38;2;243;243;255m\x1b[48;2;240;238;252m▀\x1b[0m\x1b[38;2;245;243;255m\x1b[48;2;244;242;254m▀\x1b[0m\x1b[38;2;247;246;251m\x1b[48;2;245;243;254m▀\x1b[0m\x1b[38;2;249;247;251m\x1b[48;2;248;247;255m▀\x1b[0m\x1b[38;2;143;139;153m\x1b[48;2;244;243;249m▀\x1b[0m\x1b[38;2;247;244;255m▄\x1b[0m            \x1b[38;2;227;222;244m▄\x1b[0m\x1b[38;2;229;227;240m▀\x1b[0m       \n" +
	"         \x1b[38;2;240;238;253m▀\x1b[0m\x1b[38;2;241;238;255m\x1b[48;2;239;237;251m▀\x1b[0m\x1b[38;2;244;242;254m\x1b[48;2;241;239;253m▀\x1b[0m\x1b[38;2;245;243;254m\x1b[48;2;243;241;255m▀\x1b[0m\x1b[38;2;246;245;252m\x1b[48;2;244;242;255m▀\x1b[0m\x1b[38;2;248;247;252m\x1b[48;2;245;244;253m▀\x1b[0m\x1b[38;2;249;248;254m\x1b[48;2;245;244;253m▀\x1b[0m\x1b[38;2;112;110;121m\x1b[48;2;245;244;252m▀\x1b[0m\x1b[38;2;247;246;252m▄\x1b[0m\x1b[38;2;247;246;253m▄\x1b[0m\x1b[38;2;253;251;255m▄\x1b[0m\x1b[38;2;234;232;245m▄\x1b[0m\x1b[38;2;242;240;254m▄\x1b[0m\x1b[38;2;244;242;253m▄\x1b[0m\x1b[38;2;245;243;254m▄\x1b[0m\x1b[38;2;237;234;252m▄\x1b[0m\x1b[38;2;208;206;219m\x1b[48;2;233;229;252m▀\x1b[0m\x1b[38;2;243;240;255m\x1b[48;2;93;91;105m▀\x1b[0m\x1b[38;2;88;86;99m▀\x1b[0m        \n" +
	"           \x1b[38;2;121;119;133m▀\x1b[0m\x1b[38;2;236;233;250m▀\x1b[0m\x1b[38;2;241;238;255m\x1b[48;2;165;163;177m▀\x1b[0m\x1b[38;2;242;241;255m\x1b[48;2;243;240;255m▀\x1b[0m\x1b[38;2;243;243;253m\x1b[48;2;239;236;254m▀\x1b[0m\x1b[38;2;243;243;254m\x1b[48;2;240;237;254m▀\x1b[0m\x1b[38;2;243;243;254m\x1b[48;2;241;238;255m▀\x1b[0m\x1b[38;2;243;243;255m\x1b[48;2;241;238;255m▀\x1b[0m\x1b[38;2;243;243;255m\x1b[48;2;240;237;254m▀\x1b[0m\x1b[38;2;243;241;254m\x1b[48;2;238;235;254m▀\x1b[0m\x1b[38;2;242;239;255m\x1b[48;2;239;236;253m▀\x1b[0m\x1b[38;2;240;237;254m\x1b[48;2;248;245;255m▀\x1b[0m\x1b[38;2;239;236;255m\x1b[48;2;96;94;108m▀\x1b[0m\x1b[38;2;250;248;255m▀\x1b[0m           \n" +
	"                \x1b[38;2;110;108;122m▀\x1b[0m\x1b[38;2;93;91;105m▀\x1b[0m\x1b[38;2;88;86;100m▀\x1b[0m\x1b[38;2;105;103;117m▀\x1b[0m                \n"

// iconWidth is the column count the art was rendered at. Callers need it to lay
// out beside the icon without measuring escape sequences.
const iconWidth = 36
