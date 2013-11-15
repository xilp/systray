#import <Foundation/Foundation.h> 

@interface Communicator : NSObject <NSStreamDelegate>
{
@public
    NSString *host; 
    int port; 
} 

- (void)setup; 
- (void)open; 
- (void)close; 
- (void)stream:(NSStream *)stream handleEvent:(NSStreamEvent)event; 
- (void)readIn:(NSString *)s; 
- (void)writeOut:(NSString *)s; 

@end 
