import * as cdk from 'aws-cdk-lib'
import { Construct } from 'constructs'
import * as cognito from 'aws-cdk-lib/aws-cognito'

export class CdkStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props)

    const pool = new cognito.UserPool(this, 'UserPool', {
      userPoolName: 'reservation-system',
      signInCaseSensitive: false,
      selfSignUpEnabled: true,
      userVerification: {
        emailSubject: 'Verify your email for our awesome app!',
        emailBody: 'Thanks for signing up to our awesome app! Your verification code is {####}',
        emailStyle: cognito.VerificationEmailStyle.CODE,
        smsMessage: 'Thanks for signing up to our awesome app! Your verification code is {####}',
      },
      signInAliases: {
        email: true,
      },
      autoVerify: { email: true },
      standardAttributes: {
        email: {
          required: true,
          mutable: false,
        },
      },
    })

    pool.addClient('reservation-system-frontend', {
      authFlows: {
        userPassword: true,
        userSrp: true,
      },
      generateSecret: true,
      supportedIdentityProviders: [cognito.UserPoolClientIdentityProvider.COGNITO],
      oAuth: {
        flows: {
          authorizationCodeGrant: true,
        },
        scopes: [ cognito.OAuthScope.EMAIL, cognito.OAuthScope.OPENID, cognito.OAuthScope.PROFILE ],
        callbackUrls: [ 'http://localhost:3000/api/auth/callback/cognito' ],
        logoutUrls: [ 'http://localhost:3000/' ],
      },
    })

    pool.addDomain('reservation-system-domain', {
      cognitoDomain: {
        domainPrefix: 'reservation-system',
      },
    })
  }
}
