package packed

import "github.com/gogf/gf/v2/os/gres"

func init() {
	if err := gres.Add("H4sIAAAAAAAC/wrwZmYRYeBgYGAQC9CPZkAC7AycDJmGFnn6qXmhIawMjGceH88I8GbnQFYD082BoVsQoVs/KbE4Va8kPzcHbk6Rl44fr5fuN19dPT8f/9ANoQFaer5+frp6YUEPmFgYGP7/B9n0YaaMnhIDA4MMkk2Y7uSC2VSVoevsR4ZTRVAMIMq1VdcMH2Zer1Fl1XNggLk2faJxDsi10iiunWlrhWIZC9QysOFFLtewuJORSYQZd5RAgADDW0cQjRlBMN0QfyIHoSBcNwPDkkY7hG4sXsbnCi4UV6yAmYMe/KgOQQ4dERSHvEQxgIBbsAUnzC3/HcMYGdADl5UNJMnKwMrgy8jAMJERxAMEAAD//y7SzxX0AgAA"); err != nil {
		panic("add binary content to resource manager failed: " + err.Error())
	}
}
