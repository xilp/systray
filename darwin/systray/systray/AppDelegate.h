#import <Cocoa/Cocoa.h>

@interface AppDelegate : NSObject <NSApplicationDelegate, NSStreamDelegate>
{
    NSStatusItem *statusItem;
    NSTimer *timer;
    NSMutableDictionary *icons;

    NSString *host;
    int port;
    bool unconnected;

    NSInputStream *inputStream;  
    NSOutputStream *outputStream;

    int headReading;
    uint8_t headCache[4];
    int bodyReading;
    NSMutableData *bodyCache;
}

- (IBAction)clicked:(id)sender;

- (void)open;
- (void)close;

- (void)send:(NSMutableDictionary*)cmd;
- (void)received:(NSData*)data;

- (void)stream:(NSStream*)stream handleEvent:(NSStreamEvent)event;
- (void)handleTimer:(NSTimer*)timer;

@property (assign) IBOutlet NSWindow *window;

@end
