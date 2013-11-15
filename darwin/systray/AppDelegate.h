#import <Cocoa/Cocoa.h>

@interface AppDelegate : NSObject <NSApplicationDelegate>
{
    NSStatusItem *statusItem;
}
- (IBAction)open:(id)sender;

@property (assign) IBOutlet NSWindow *window;

@end
