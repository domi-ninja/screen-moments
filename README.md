# Screen Moments

Sometimes, I play video games with friends. Sometimes, this is really funny. Rarely, I like to cut together some clips of the best moments using the excellent [losslesscut](https://mifi.no/losslesscut/) video cutting software and show it to them in our chat.

I did not realize steam has a [rolling buffer screen recording feature now](https://help.steampowered.com/en/faqs/view/23B7-49AD-4A28-9590) (which is very cool and fills this niche nicely). 

Because I did not know about the steam feature, I came up with a plan to implement this on my own:

1. Use OBS to record (because I sure did not feel like implementing screen recording and video compression **well** myself)
2. Have OBS target localhost as a streaming target and keep that rolling the entire time
3. Write a program that recieves this stream and listens to hotkeys to permanently keep the last X minutes of stream on disk

I have had a suprising amount of success having claude vibe code this entire thing in golang. It is not 100% finished, and I am bailing out of this project for reasons above, but I am uploading the result for future reference, in case I need to make some other small utility program for windows.


## how to run 

`air`

## how to target

1. Open OBS Studio
2. Go to Settings -> Stream
3. Select "Custom..." as the service
4. Set the Server to: `rtmp://localhost:1935/live`
5. Set the Stream Key to: `screenmoments`
6. Click "OK" to save settings
7. Click "Start Streaming" in OBS
