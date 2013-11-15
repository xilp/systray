#import "AppDelegate.h"
#import "Communicator.h"

@implementation AppDelegate

@synthesize window = _window;

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification
{
    statusItem = [[NSStatusBar systemStatusBar] statusItemWithLength:NSVariableStatusItemLength];
    [statusItem setTitle:@"D"];
    [statusItem setAction:@selector(open:)];
    [statusItem setHighlightMode:YES];
}

- (IBAction)open:(id)sender {
    NSURL *url = [NSURL URLWithString:@"http://mail.google.com"];
    [[NSWorkspace sharedWorkspace] openURL:url];
    NSLog(@"click!");
}

@end
