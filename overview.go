//
// openzwave provides a thin Go wrapper around the openzwave library.
//
// The following shows a simple use of the API which will log every notification received.
//
//        var loop = func(api openzwave.API) {
//                fmt.Printf("event loop starts\n")
//                for {
//                        select {
//                        case quitNow := <-api.QuitSignal():
//                                _ = quitNow
//                                fmt.Printf("event loop ends\n")
//                                return
//                        }
//                }
//        }
//
//        os.Exit(openzwave.
//                BuildAPI("../go-openzwave/openzwave/config", "", "").
//                AddIntOption("SaveLogLevel", LOG_LEVEL.NONE).
//                AddIntOption("QueueLogLevel", LOG_LEVEL.NONE).
//                AddIntOption("DumpTrigger", LOG_LEVEL.NONE).
//                AddIntOption("PollInterval", 500).
//                AddBoolOption("IntervalBetweenPolls", true).
//                AddBoolOption("ValidateValueChanges", true).
//                Run(loop))
package openzwave
