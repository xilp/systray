#import <Cocoa/Cocoa.h>

@interface AppDelegate : NSObject <NSApplicationDelegate, NSStreamDelegate>
{
    NSStatusItem *statusItem;

    NSString *host;
    int port;

    NSInputStream *inputStream;  
    NSOutputStream *outputStream;

    int headReading;
    uint8_t headCache[4];
    int bodyReading;
    NSMutableData* bodyCache;
}

- (IBAction)clicked:(id)sender;

- (void)open;
- (void)close;

- (void)received:(NSString*)data;

- (void)stream:(NSStream*)stream handleEvent:(NSStreamEvent)event;

@property (assign) IBOutlet NSWindow *window;

@end
