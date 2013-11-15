#import <Foundation/Foundation.h> 

@interface Communicator : NSObject <NSStreamDelegate>
{
    NSString *host;
    int port;
}

- (void)setup:(NSString *)host port:(int)port;
- (void)open;
- (void)close;
- (void)stream:(NSStream *)stream event:(NSStreamEvent)event;

@end 
