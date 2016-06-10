package aws_test

import (
	"os"
	"testing"
	"time"

	"github.com/convox/rack/api/awsutil"
	"github.com/convox/rack/api/provider"
	"github.com/convox/rack/api/structs"
	"github.com/stretchr/testify/assert"
)

func init() {
	os.Setenv("CLUSTER", "convox-0-Cluster-QMIL82M9O5TK")
	os.Setenv("RACK", "convox-0")
}

func TestInstanceListAPI(t *testing.T) {
	defer func() {
		provider.CurrentProvider = new(provider.TestProviderRunner)
	}()

	aws := StubAwsProvider(
		describeInstances,          // http://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeInstances.html
		listContainerInstances,     // http://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_ListContainerInstances.html
		describeContainerInstances, // http://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_DescribeContainerInstances.html
		// describeAutoScalingGroups,  // http://docs.aws.amazon.com/AutoScaling/latest/APIReference/API_DescribeAutoScalingGroups.html
	)
	defer aws.Close()

	is, err := provider.InstanceList()
	assert.Nil(t, err)
	assert.EqualValues(t,
		structs.Instances{
			structs.Instance{
				Agent:     false,
				Cpu:       0.0,
				Id:        "i-3a64b9c2",
				Memory:    0.0,
				PrivateIp: "10.0.3.123",
				Processes: 2,
				PublicIp:  "54.175.48.124",
				Status:    "",
				Started:   time.Unix(1465521109, 0).UTC(),
			},
			structs.Instance{
				Agent:     false,
				Cpu:       0.0,
				Id:        "i-a8db21ed",
				Memory:    0.0,
				PrivateIp: "10.0.1.158",
				Processes: 1,
				PublicIp:  "107.21.88.144",
				Status:    "",
				Started:   time.Unix(1465521106, 0).UTC(),
			},
			structs.Instance{
				Agent:     false,
				Cpu:       0.0,
				Id:        "i-fac3ba60",
				Memory:    0.0,
				PrivateIp: "10.0.2.82",
				Processes: 0,
				PublicIp:  "52.90.183.103",
				Status:    "",
				Started:   time.Unix(1465521110, 0).UTC(),
			},
		},
		is,
	)
}

func TestInstanceList(t *testing.T) {
	defer func() {
		provider.CurrentProvider = new(provider.TestProviderRunner)
	}()

	p := provider.CurrentProvider.(*provider.TestProviderRunner)
	p.Mock.On("InstanceList").Return(
		structs.Instances{
			structs.Instance{
				Agent:     false,
				Cpu:       0.0,
				Id:        "i-3a64b9c2",
				Memory:    0.0,
				PrivateIp: "",
				Processes: 2,
				PublicIp:  "",
				Status:    "",
				Started:   time.Unix(1465521109, 0).UTC(),
			},
			structs.Instance{
				Agent:     false,
				Cpu:       0.0,
				Id:        "i-a8db21ed",
				Memory:    0.0,
				PrivateIp: "",
				Processes: 1,
				PublicIp:  "",
				Status:    "",
				Started:   time.Unix(1465521106, 0).UTC(),
			},
			structs.Instance{
				Agent:     false,
				Cpu:       0.0,
				Id:        "i-fac3ba60",
				Memory:    0.0,
				PrivateIp: "",
				Processes: 0,
				PublicIp:  "",
				Status:    "",
				Started:   time.Unix(1465521110, 0).UTC(),
			},
		},
		nil)

	is, err := provider.InstanceList()
	assert.Nil(t, err)
	assert.Equal(t, 3, len(is))
}

/*
AWS Request / Response Cycles

# Use CI account
$ export AWS_DEFAULT_PROFILE=ci
$ AWS_ACCESS_KEY_ID=$(aws configure get aws_access_key_id)
$ AWS_SECRET_ACCESS_KEY=$(aws configure get aws_secret_access_key)

$ convox install --stack-name=convox-0
Installing Convox (20160607215357)...
...
Created ECS Cluster: convox-0-Cluster-QMIL82M9O5TK
Created AutoScalingGroup: convox-0-Instances-1NDOC550I1C8H
*/

// $ aws --debug ec2 describe-instances --filters "Name=tag:Name,Values=convox-0"
var describeInstances = awsutil.Cycle{
	Request: awsutil.Request{
		RequestURI: "/",
		Operation:  "",
		Body:       `Action=DescribeInstances&Filter.1.Name=tag&Filter.1.Value.1=convox-0&Version=2015-10-01`,
	},
	Response: awsutil.Response{
		StatusCode: 200,
		Body: `<DescribeInstancesResponse xmlns="http://ec2.amazonaws.com/doc/2015-10-01/">
    <requestId>38b81354-aa93-433c-9b0b-dc4388184897</requestId>
    <reservationSet>
        <item>
            <reservationId>r-0dc796d8</reservationId>
            <ownerId>568149725493</ownerId>
            <groupSet/>
            <instancesSet>
                <item>
                    <instanceId>i-a8db21ed</instanceId>
                    <imageId>ami-67a3a90d</imageId>
                    <instanceState>
                        <code>16</code>
                        <name>running</name>
                    </instanceState>
                    <privateDnsName>ip-10-0-1-158.ec2.internal</privateDnsName>
                    <dnsName>ec2-107-21-88-144.compute-1.amazonaws.com</dnsName>
                    <reason/>
                    <amiLaunchIndex>0</amiLaunchIndex>
                    <productCodes/>
                    <instanceType>t2.small</instanceType>
                    <launchTime>2016-06-10T01:11:46.000Z</launchTime>
                    <placement>
                        <availabilityZone>us-east-1a</availabilityZone>
                        <groupName/>
                        <tenancy>default</tenancy>
                    </placement>
                    <monitoring>
                        <state>enabled</state>
                    </monitoring>
                    <subnetId>subnet-092e127f</subnetId>
                    <vpcId>vpc-0ca2056b</vpcId>
                    <privateIpAddress>10.0.1.158</privateIpAddress>
                    <ipAddress>107.21.88.144</ipAddress>
                    <sourceDestCheck>true</sourceDestCheck>
                    <groupSet>
                        <item>
                            <groupId>sg-7a990801</groupId>
                            <groupName>convox-0-SecurityGroup-1FNBUJMAS7PDR</groupName>
                        </item>
                    </groupSet>
                    <architecture>x86_64</architecture>
                    <rootDeviceType>ebs</rootDeviceType>
                    <rootDeviceName>/dev/xvda</rootDeviceName>
                    <blockDeviceMapping>
                        <item>
                            <deviceName>/dev/xvda</deviceName>
                            <ebs>
                                <volumeId>vol-5a5bacf0</volumeId>
                                <status>attached</status>
                                <attachTime>2016-06-10T01:11:47.000Z</attachTime>
                                <deleteOnTermination>true</deleteOnTermination>
                            </ebs>
                        </item>
                        <item>
                            <deviceName>/dev/sdb</deviceName>
                            <ebs>
                                <volumeId>vol-b85aad12</volumeId>
                                <status>attached</status>
                                <attachTime>2016-06-10T01:11:47.000Z</attachTime>
                                <deleteOnTermination>true</deleteOnTermination>
                            </ebs>
                        </item>
                        <item>
                            <deviceName>/dev/xvdcz</deviceName>
                            <ebs>
                                <volumeId>vol-1e5bacb4</volumeId>
                                <status>attached</status>
                                <attachTime>2016-06-10T01:11:47.000Z</attachTime>
                                <deleteOnTermination>true</deleteOnTermination>
                            </ebs>
                        </item>
                    </blockDeviceMapping>
                    <virtualizationType>hvm</virtualizationType>
                    <clientToken>7d24094b-a575-4d92-ae0d-91953f407266_subnet-092e127f_1</clientToken>
                    <tagSet>
                        <item>
                            <key>Name</key>
                            <value>convox-0</value>
                        </item>
                        <item>
                            <key>aws:cloudformation:stack-id</key>
                            <value>arn:aws:cloudformation:us-east-1:568149725493:stack/convox-0/bdf699d0-2ea7-11e6-9689-500c28903236</value>
                        </item>
                        <item>
                            <key>Rack</key>
                            <value>convox-0</value>
                        </item>
                        <item>
                            <key>aws:cloudformation:logical-id</key>
                            <value>Instances</value>
                        </item>
                        <item>
                            <key>aws:cloudformation:stack-name</key>
                            <value>convox-0</value>
                        </item>
                        <item>
                            <key>aws:autoscaling:groupName</key>
                            <value>convox-0-Instances-1NDOC550I1C8H</value>
                        </item>
                    </tagSet>
                    <hypervisor>xen</hypervisor>
                    <networkInterfaceSet>
                        <item>
                            <networkInterfaceId>eni-ea0feaa8</networkInterfaceId>
                            <subnetId>subnet-092e127f</subnetId>
                            <vpcId>vpc-0ca2056b</vpcId>
                            <description/>
                            <ownerId>568149725493</ownerId>
                            <status>in-use</status>
                            <macAddress>0a:df:7f:e4:26:85</macAddress>
                            <privateIpAddress>10.0.1.158</privateIpAddress>
                            <privateDnsName>ip-10-0-1-158.ec2.internal</privateDnsName>
                            <sourceDestCheck>true</sourceDestCheck>
                            <groupSet>
                                <item>
                                    <groupId>sg-7a990801</groupId>
                                    <groupName>convox-0-SecurityGroup-1FNBUJMAS7PDR</groupName>
                                </item>
                            </groupSet>
                            <attachment>
                                <attachmentId>eni-attach-6275bfb3</attachmentId>
                                <deviceIndex>0</deviceIndex>
                                <status>attached</status>
                                <attachTime>2016-06-10T01:11:46.000Z</attachTime>
                                <deleteOnTermination>true</deleteOnTermination>
                            </attachment>
                            <association>
                                <publicIp>107.21.88.144</publicIp>
                                <publicDnsName>ec2-107-21-88-144.compute-1.amazonaws.com</publicDnsName>
                                <ipOwnerId>amazon</ipOwnerId>
                            </association>
                            <privateIpAddressesSet>
                                <item>
                                    <privateIpAddress>10.0.1.158</privateIpAddress>
                                    <privateDnsName>ip-10-0-1-158.ec2.internal</privateDnsName>
                                    <primary>true</primary>
                                    <association>
                                    <publicIp>107.21.88.144</publicIp>
                                    <publicDnsName>ec2-107-21-88-144.compute-1.amazonaws.com</publicDnsName>
                                    <ipOwnerId>amazon</ipOwnerId>
                                    </association>
                                </item>
                            </privateIpAddressesSet>
                        </item>
                    </networkInterfaceSet>
                    <iamInstanceProfile>
                        <arn>arn:aws:iam::568149725493:instance-profile/convox-0-InstanceProfile-66WMQ9ESABCM</arn>
                        <id>AIPAIKSK6M654S3P5CW7M</id>
                    </iamInstanceProfile>
                    <ebsOptimized>false</ebsOptimized>
                </item>
            </instancesSet>
            <requesterId>226008221399</requesterId>
        </item>
        <item>
            <reservationId>r-efd79816</reservationId>
            <ownerId>568149725493</ownerId>
            <groupSet/>
            <instancesSet>
                <item>
                    <instanceId>i-3a64b9c2</instanceId>
                    <imageId>ami-67a3a90d</imageId>
                    <instanceState>
                        <code>16</code>
                        <name>running</name>
                    </instanceState>
                    <privateDnsName>ip-10-0-3-123.ec2.internal</privateDnsName>
                    <dnsName>ec2-54-175-48-124.compute-1.amazonaws.com</dnsName>
                    <reason/>
                    <amiLaunchIndex>0</amiLaunchIndex>
                    <productCodes/>
                    <instanceType>t2.small</instanceType>
                    <launchTime>2016-06-10T01:11:49.000Z</launchTime>
                    <placement>
                        <availabilityZone>us-east-1e</availabilityZone>
                        <groupName/>
                        <tenancy>default</tenancy>
                    </placement>
                    <monitoring>
                        <state>enabled</state>
                    </monitoring>
                    <subnetId>subnet-9e5f75a3</subnetId>
                    <vpcId>vpc-0ca2056b</vpcId>
                    <privateIpAddress>10.0.3.123</privateIpAddress>
                    <ipAddress>54.175.48.124</ipAddress>
                    <sourceDestCheck>true</sourceDestCheck>
                    <groupSet>
                        <item>
                            <groupId>sg-7a990801</groupId>
                            <groupName>convox-0-SecurityGroup-1FNBUJMAS7PDR</groupName>
                        </item>
                    </groupSet>
                    <architecture>x86_64</architecture>
                    <rootDeviceType>ebs</rootDeviceType>
                    <rootDeviceName>/dev/xvda</rootDeviceName>
                    <blockDeviceMapping>
                        <item>
                            <deviceName>/dev/xvda</deviceName>
                            <ebs>
                                <volumeId>vol-63dfeeed</volumeId>
                                <status>attached</status>
                                <attachTime>2016-06-10T01:11:50.000Z</attachTime>
                                <deleteOnTermination>true</deleteOnTermination>
                            </ebs>
                        </item>
                        <item>
                            <deviceName>/dev/sdb</deviceName>
                            <ebs>
                                <volumeId>vol-09dfee87</volumeId>
                                <status>attached</status>
                                <attachTime>2016-06-10T01:11:50.000Z</attachTime>
                                <deleteOnTermination>true</deleteOnTermination>
                            </ebs>
                        </item>
                        <item>
                            <deviceName>/dev/xvdcz</deviceName>
                            <ebs>
                                <volumeId>vol-24dfeeaa</volumeId>
                                <status>attached</status>
                                <attachTime>2016-06-10T01:11:50.000Z</attachTime>
                                <deleteOnTermination>true</deleteOnTermination>
                            </ebs>
                        </item>
                    </blockDeviceMapping>
                    <virtualizationType>hvm</virtualizationType>
                    <clientToken>7d24094b-a575-4d92-ae0d-91953f407266_subnet-9e5f75a3_1</clientToken>
                    <tagSet>
                        <item>
                            <key>aws:cloudformation:stack-id</key>
                            <value>arn:aws:cloudformation:us-east-1:568149725493:stack/convox-0/bdf699d0-2ea7-11e6-9689-500c28903236</value>
                        </item>
                        <item>
                            <key>aws:autoscaling:groupName</key>
                            <value>convox-0-Instances-1NDOC550I1C8H</value>
                        </item>
                        <item>
                            <key>Name</key>
                            <value>convox-0</value>
                        </item>
                        <item>
                            <key>aws:cloudformation:stack-name</key>
                            <value>convox-0</value>
                        </item>
                        <item>
                            <key>Rack</key>
                            <value>convox-0</value>
                        </item>
                        <item>
                            <key>aws:cloudformation:logical-id</key>
                            <value>Instances</value>
                        </item>
                    </tagSet>
                    <hypervisor>xen</hypervisor>
                    <networkInterfaceSet>
                        <item>
                            <networkInterfaceId>eni-1faa7e21</networkInterfaceId>
                            <subnetId>subnet-9e5f75a3</subnetId>
                            <vpcId>vpc-0ca2056b</vpcId>
                            <description/>
                            <ownerId>568149725493</ownerId>
                            <status>in-use</status>
                            <macAddress>06:2f:ec:39:90:bb</macAddress>
                            <privateIpAddress>10.0.3.123</privateIpAddress>
                            <privateDnsName>ip-10-0-3-123.ec2.internal</privateDnsName>
                            <sourceDestCheck>true</sourceDestCheck>
                            <groupSet>
                                <item>
                                    <groupId>sg-7a990801</groupId>
                                    <groupName>convox-0-SecurityGroup-1FNBUJMAS7PDR</groupName>
                                </item>
                            </groupSet>
                            <attachment>
                                <attachmentId>eni-attach-38de8ad1</attachmentId>
                                <deviceIndex>0</deviceIndex>
                                <status>attached</status>
                                <attachTime>2016-06-10T01:11:49.000Z</attachTime>
                                <deleteOnTermination>true</deleteOnTermination>
                            </attachment>
                            <association>
                                <publicIp>54.175.48.124</publicIp>
                                <publicDnsName>ec2-54-175-48-124.compute-1.amazonaws.com</publicDnsName>
                                <ipOwnerId>amazon</ipOwnerId>
                            </association>
                            <privateIpAddressesSet>
                                <item>
                                    <privateIpAddress>10.0.3.123</privateIpAddress>
                                    <privateDnsName>ip-10-0-3-123.ec2.internal</privateDnsName>
                                    <primary>true</primary>
                                    <association>
                                    <publicIp>54.175.48.124</publicIp>
                                    <publicDnsName>ec2-54-175-48-124.compute-1.amazonaws.com</publicDnsName>
                                    <ipOwnerId>amazon</ipOwnerId>
                                    </association>
                                </item>
                            </privateIpAddressesSet>
                        </item>
                    </networkInterfaceSet>
                    <iamInstanceProfile>
                        <arn>arn:aws:iam::568149725493:instance-profile/convox-0-InstanceProfile-66WMQ9ESABCM</arn>
                        <id>AIPAIKSK6M654S3P5CW7M</id>
                    </iamInstanceProfile>
                    <ebsOptimized>false</ebsOptimized>
                </item>
            </instancesSet>
            <requesterId>226008221399</requesterId>
        </item>
        <item>
            <reservationId>r-1e15e6bd</reservationId>
            <ownerId>568149725493</ownerId>
            <groupSet/>
            <instancesSet>
                <item>
                    <instanceId>i-fac3ba60</instanceId>
                    <imageId>ami-67a3a90d</imageId>
                    <instanceState>
                        <code>16</code>
                        <name>running</name>
                    </instanceState>
                    <privateDnsName>ip-10-0-2-82.ec2.internal</privateDnsName>
                    <dnsName>ec2-52-90-183-103.compute-1.amazonaws.com</dnsName>
                    <reason/>
                    <amiLaunchIndex>0</amiLaunchIndex>
                    <productCodes/>
                    <instanceType>t2.small</instanceType>
                    <launchTime>2016-06-10T01:11:50.000Z</launchTime>
                    <placement>
                        <availabilityZone>us-east-1b</availabilityZone>
                        <groupName/>
                        <tenancy>default</tenancy>
                    </placement>
                    <monitoring>
                        <state>enabled</state>
                    </monitoring>
                    <subnetId>subnet-ff89c0a7</subnetId>
                    <vpcId>vpc-0ca2056b</vpcId>
                    <privateIpAddress>10.0.2.82</privateIpAddress>
                    <ipAddress>52.90.183.103</ipAddress>
                    <sourceDestCheck>true</sourceDestCheck>
                    <groupSet>
                        <item>
                            <groupId>sg-7a990801</groupId>
                            <groupName>convox-0-SecurityGroup-1FNBUJMAS7PDR</groupName>
                        </item>
                    </groupSet>
                    <architecture>x86_64</architecture>
                    <rootDeviceType>ebs</rootDeviceType>
                    <rootDeviceName>/dev/xvda</rootDeviceName>
                    <blockDeviceMapping>
                        <item>
                            <deviceName>/dev/xvda</deviceName>
                            <ebs>
                                <volumeId>vol-4eb8e1eb</volumeId>
                                <status>attached</status>
                                <attachTime>2016-06-10T01:11:51.000Z</attachTime>
                                <deleteOnTermination>true</deleteOnTermination>
                            </ebs>
                        </item>
                        <item>
                            <deviceName>/dev/sdb</deviceName>
                            <ebs>
                                <volumeId>vol-efbbe24a</volumeId>
                                <status>attached</status>
                                <attachTime>2016-06-10T01:11:51.000Z</attachTime>
                                <deleteOnTermination>true</deleteOnTermination>
                            </ebs>
                        </item>
                        <item>
                            <deviceName>/dev/xvdcz</deviceName>
                            <ebs>
                                <volumeId>vol-a1bbe204</volumeId>
                                <status>attached</status>
                                <attachTime>2016-06-10T01:11:51.000Z</attachTime>
                                <deleteOnTermination>true</deleteOnTermination>
                            </ebs>
                        </item>
                    </blockDeviceMapping>
                    <virtualizationType>hvm</virtualizationType>
                    <clientToken>7d24094b-a575-4d92-ae0d-91953f407266_subnet-ff89c0a7_1</clientToken>
                    <tagSet>
                        <item>
                            <key>aws:cloudformation:logical-id</key>
                            <value>Instances</value>
                        </item>
                        <item>
                            <key>aws:cloudformation:stack-name</key>
                            <value>convox-0</value>
                        </item>
                        <item>
                            <key>Name</key>
                            <value>convox-0</value>
                        </item>
                        <item>
                            <key>Rack</key>
                            <value>convox-0</value>
                        </item>
                        <item>
                            <key>aws:autoscaling:groupName</key>
                            <value>convox-0-Instances-1NDOC550I1C8H</value>
                        </item>
                        <item>
                            <key>aws:cloudformation:stack-id</key>
                            <value>arn:aws:cloudformation:us-east-1:568149725493:stack/convox-0/bdf699d0-2ea7-11e6-9689-500c28903236</value>
                        </item>
                    </tagSet>
                    <hypervisor>xen</hypervisor>
                    <networkInterfaceSet>
                        <item>
                            <networkInterfaceId>eni-b23f95e0</networkInterfaceId>
                            <subnetId>subnet-ff89c0a7</subnetId>
                            <vpcId>vpc-0ca2056b</vpcId>
                            <description/>
                            <ownerId>568149725493</ownerId>
                            <status>in-use</status>
                            <macAddress>0e:ab:6b:1d:45:b1</macAddress>
                            <privateIpAddress>10.0.2.82</privateIpAddress>
                            <privateDnsName>ip-10-0-2-82.ec2.internal</privateDnsName>
                            <sourceDestCheck>true</sourceDestCheck>
                            <groupSet>
                                <item>
                                    <groupId>sg-7a990801</groupId>
                                    <groupName>convox-0-SecurityGroup-1FNBUJMAS7PDR</groupName>
                                </item>
                            </groupSet>
                            <attachment>
                                <attachmentId>eni-attach-d607a828</attachmentId>
                                <deviceIndex>0</deviceIndex>
                                <status>attached</status>
                                <attachTime>2016-06-10T01:11:50.000Z</attachTime>
                                <deleteOnTermination>true</deleteOnTermination>
                            </attachment>
                            <association>
                                <publicIp>52.90.183.103</publicIp>
                                <publicDnsName>ec2-52-90-183-103.compute-1.amazonaws.com</publicDnsName>
                                <ipOwnerId>amazon</ipOwnerId>
                            </association>
                            <privateIpAddressesSet>
                                <item>
                                    <privateIpAddress>10.0.2.82</privateIpAddress>
                                    <privateDnsName>ip-10-0-2-82.ec2.internal</privateDnsName>
                                    <primary>true</primary>
                                    <association>
                                    <publicIp>52.90.183.103</publicIp>
                                    <publicDnsName>ec2-52-90-183-103.compute-1.amazonaws.com</publicDnsName>
                                    <ipOwnerId>amazon</ipOwnerId>
                                    </association>
                                </item>
                            </privateIpAddressesSet>
                        </item>
                    </networkInterfaceSet>
                    <iamInstanceProfile>
                        <arn>arn:aws:iam::568149725493:instance-profile/convox-0-InstanceProfile-66WMQ9ESABCM</arn>
                        <id>AIPAIKSK6M654S3P5CW7M</id>
                    </iamInstanceProfile>
                    <ebsOptimized>false</ebsOptimized>
                </item>
            </instancesSet>
            <requesterId>226008221399</requesterId>
        </item>
    </reservationSet>
</DescribeInstancesResponse>`},
}

// $ aws --debug ecs list-container-instances --cluster=convox-0-Cluster-QMIL82M9O5TK
var listContainerInstances = awsutil.Cycle{
	Request: awsutil.Request{
		RequestURI: "/",
		Operation:  "AmazonEC2ContainerServiceV20141113.ListContainerInstances",
		Body:       `{"cluster": "convox-0-Cluster-QMIL82M9O5TK"}`,
	},
	Response: awsutil.Response{
		StatusCode: 200,
		Body: `{
    "containerInstanceArns": [
        "arn:aws:ecs:us-east-1:568149725493:container-instance/2902d068-fa56-45b3-b6d0-4f19b39d02b7",
        "arn:aws:ecs:us-east-1:568149725493:container-instance/848937b6-e0a0-4036-babe-206a63332838",
        "arn:aws:ecs:us-east-1:568149725493:container-instance/c7ec7762-ecc7-4bea-b44d-ba64b8ce547d"
    ]
}`},
}

// $ aws --debug ecs describe-container-instances \
// --cluster=convox-0-Cluster-QMIL82M9O5TK \
// --container-instances \
//   arn:aws:ecs:us-east-1:568149725493:container-instance/2902d068-fa56-45b3-b6d0-4f19b39d02b7 \
//   arn:aws:ecs:us-east-1:568149725493:container-instance/848937b6-e0a0-4036-babe-206a63332838 \
//   arn:aws:ecs:us-east-1:568149725493:container-instance/c7ec7762-ecc7-4bea-b44d-ba64b8ce547d
var describeContainerInstances = awsutil.Cycle{
	Request: awsutil.Request{
		RequestURI: "/",
		Operation:  "AmazonEC2ContainerServiceV20141113.DescribeContainerInstances",
		Body:       `{"cluster": "convox-0-Cluster-QMIL82M9O5TK", "containerInstances": ["arn:aws:ecs:us-east-1:568149725493:container-instance/2902d068-fa56-45b3-b6d0-4f19b39d02b7", "arn:aws:ecs:us-east-1:568149725493:container-instance/848937b6-e0a0-4036-babe-206a63332838", "arn:aws:ecs:us-east-1:568149725493:container-instance/c7ec7762-ecc7-4bea-b44d-ba64b8ce547d"]}`,
	},
	Response: awsutil.Response{
		StatusCode: 200,
		Body: `{
    "failures": [],
    "containerInstances": [
        {
            "status": "ACTIVE",
            "registeredResources": [
                {
                    "integerValue": 1024,
                    "longValue": 0,
                    "type": "INTEGER",
                    "name": "CPU",
                    "doubleValue": 0.0
                },
                {
                    "integerValue": 2003,
                    "longValue": 0,
                    "type": "INTEGER",
                    "name": "MEMORY",
                    "doubleValue": 0.0
                },
                {
                    "name": "PORTS",
                    "longValue": 0,
                    "doubleValue": 0.0,
                    "stringSetValue": [
                        "22",
                        "2376",
                        "2375",
                        "51678"
                    ],
                    "type": "STRINGSET",
                    "integerValue": 0
                },
                {
                    "name": "PORTS_UDP",
                    "longValue": 0,
                    "doubleValue": 0.0,
                    "stringSetValue": [],
                    "type": "STRINGSET",
                    "integerValue": 0
                }
            ],
            "ec2InstanceId": "i-3a64b9c2",
            "agentConnected": true,
            "containerInstanceArn": "arn:aws:ecs:us-east-1:568149725493:container-instance/2902d068-fa56-45b3-b6d0-4f19b39d02b7",
            "pendingTasksCount": 0,
            "remainingResources": [
                {
                    "integerValue": 1024,
                    "longValue": 0,
                    "type": "INTEGER",
                    "name": "CPU",
                    "doubleValue": 0.0
                },
                {
                    "integerValue": 1683,
                    "longValue": 0,
                    "type": "INTEGER",
                    "name": "MEMORY",
                    "doubleValue": 0.0
                },
                {
                    "name": "PORTS",
                    "longValue": 0,
                    "doubleValue": 0.0,
                    "stringSetValue": [
                        "22",
                        "2376",
                        "2375",
                        "3101",
                        "3001",
                        "3100",
                        "51678",
                        "3000"
                    ],
                    "type": "STRINGSET",
                    "integerValue": 0
                },
                {
                    "name": "PORTS_UDP",
                    "longValue": 0,
                    "doubleValue": 0.0,
                    "stringSetValue": [],
                    "type": "STRINGSET",
                    "integerValue": 0
                }
            ],
            "runningTasksCount": 2,
            "attributes": [
                {
                    "name": "com.amazonaws.ecs.capability.privileged-container"
                },
                {
                    "name": "com.amazonaws.ecs.capability.docker-remote-api.1.17"
                },
                {
                    "name": "com.amazonaws.ecs.capability.docker-remote-api.1.18"
                },
                {
                    "name": "com.amazonaws.ecs.capability.docker-remote-api.1.19"
                },
                {
                    "name": "com.amazonaws.ecs.capability.docker-remote-api.1.20"
                },
                {
                    "name": "com.amazonaws.ecs.capability.logging-driver.json-file"
                },
                {
                    "name": "com.amazonaws.ecs.capability.logging-driver.syslog"
                },
                {
                    "name": "com.amazonaws.ecs.capability.ecr-auth"
                }
            ],
            "versionInfo": {
                "agentVersion": "1.8.2",
                "agentHash": "64205ac",
                "dockerVersion": "DockerVersion: 1.9.1"
            }
        },
        {
            "status": "ACTIVE",
            "registeredResources": [
                {
                    "integerValue": 1024,
                    "longValue": 0,
                    "type": "INTEGER",
                    "name": "CPU",
                    "doubleValue": 0.0
                },
                {
                    "integerValue": 2003,
                    "longValue": 0,
                    "type": "INTEGER",
                    "name": "MEMORY",
                    "doubleValue": 0.0
                },
                {
                    "name": "PORTS",
                    "longValue": 0,
                    "doubleValue": 0.0,
                    "stringSetValue": [
                        "22",
                        "2376",
                        "2375",
                        "51678"
                    ],
                    "type": "STRINGSET",
                    "integerValue": 0
                },
                {
                    "name": "PORTS_UDP",
                    "longValue": 0,
                    "doubleValue": 0.0,
                    "stringSetValue": [],
                    "type": "STRINGSET",
                    "integerValue": 0
                }
            ],
            "ec2InstanceId": "i-a8db21ed",
            "agentConnected": true,
            "containerInstanceArn": "arn:aws:ecs:us-east-1:568149725493:container-instance/848937b6-e0a0-4036-babe-206a63332838",
            "pendingTasksCount": 0,
            "remainingResources": [
                {
                    "integerValue": 1024,
                    "longValue": 0,
                    "type": "INTEGER",
                    "name": "CPU",
                    "doubleValue": 0.0
                },
                {
                    "integerValue": 1747,
                    "longValue": 0,
                    "type": "INTEGER",
                    "name": "MEMORY",
                    "doubleValue": 0.0
                },
                {
                    "name": "PORTS",
                    "longValue": 0,
                    "doubleValue": 0.0,
                    "stringSetValue": [
                        "22",
                        "2376",
                        "2375",
                        "3101",
                        "3001",
                        "3100",
                        "51678",
                        "3000"
                    ],
                    "type": "STRINGSET",
                    "integerValue": 0
                },
                {
                    "name": "PORTS_UDP",
                    "longValue": 0,
                    "doubleValue": 0.0,
                    "stringSetValue": [],
                    "type": "STRINGSET",
                    "integerValue": 0
                }
            ],
            "runningTasksCount": 1,
            "attributes": [
                {
                    "name": "com.amazonaws.ecs.capability.privileged-container"
                },
                {
                    "name": "com.amazonaws.ecs.capability.docker-remote-api.1.17"
                },
                {
                    "name": "com.amazonaws.ecs.capability.docker-remote-api.1.18"
                },
                {
                    "name": "com.amazonaws.ecs.capability.docker-remote-api.1.19"
                },
                {
                    "name": "com.amazonaws.ecs.capability.docker-remote-api.1.20"
                },
                {
                    "name": "com.amazonaws.ecs.capability.logging-driver.json-file"
                },
                {
                    "name": "com.amazonaws.ecs.capability.logging-driver.syslog"
                },
                {
                    "name": "com.amazonaws.ecs.capability.ecr-auth"
                }
            ],
            "versionInfo": {
                "agentVersion": "1.8.2",
                "agentHash": "64205ac",
                "dockerVersion": "DockerVersion: 1.9.1"
            }
        },
        {
            "status": "ACTIVE",
            "registeredResources": [
                {
                    "integerValue": 1024,
                    "longValue": 0,
                    "type": "INTEGER",
                    "name": "CPU",
                    "doubleValue": 0.0
                },
                {
                    "integerValue": 2003,
                    "longValue": 0,
                    "type": "INTEGER",
                    "name": "MEMORY",
                    "doubleValue": 0.0
                },
                {
                    "name": "PORTS",
                    "longValue": 0,
                    "doubleValue": 0.0,
                    "stringSetValue": [
                        "22",
                        "2376",
                        "2375",
                        "51678"
                    ],
                    "type": "STRINGSET",
                    "integerValue": 0
                },
                {
                    "name": "PORTS_UDP",
                    "longValue": 0,
                    "doubleValue": 0.0,
                    "stringSetValue": [],
                    "type": "STRINGSET",
                    "integerValue": 0
                }
            ],
            "ec2InstanceId": "i-fac3ba60",
            "agentConnected": true,
            "containerInstanceArn": "arn:aws:ecs:us-east-1:568149725493:container-instance/c7ec7762-ecc7-4bea-b44d-ba64b8ce547d",
            "pendingTasksCount": 0,
            "remainingResources": [
                {
                    "integerValue": 1024,
                    "longValue": 0,
                    "type": "INTEGER",
                    "name": "CPU",
                    "doubleValue": 0.0
                },
                {
                    "integerValue": 2003,
                    "longValue": 0,
                    "type": "INTEGER",
                    "name": "MEMORY",
                    "doubleValue": 0.0
                },
                {
                    "name": "PORTS",
                    "longValue": 0,
                    "doubleValue": 0.0,
                    "stringSetValue": [
                        "22",
                        "2376",
                        "2375",
                        "51678"
                    ],
                    "type": "STRINGSET",
                    "integerValue": 0
                },
                {
                    "name": "PORTS_UDP",
                    "longValue": 0,
                    "doubleValue": 0.0,
                    "stringSetValue": [],
                    "type": "STRINGSET",
                    "integerValue": 0
                }
            ],
            "runningTasksCount": 0,
            "attributes": [
                {
                    "name": "com.amazonaws.ecs.capability.privileged-container"
                },
                {
                    "name": "com.amazonaws.ecs.capability.docker-remote-api.1.17"
                },
                {
                    "name": "com.amazonaws.ecs.capability.docker-remote-api.1.18"
                },
                {
                    "name": "com.amazonaws.ecs.capability.docker-remote-api.1.19"
                },
                {
                    "name": "com.amazonaws.ecs.capability.docker-remote-api.1.20"
                },
                {
                    "name": "com.amazonaws.ecs.capability.logging-driver.json-file"
                },
                {
                    "name": "com.amazonaws.ecs.capability.logging-driver.syslog"
                },
                {
                    "name": "com.amazonaws.ecs.capability.ecr-auth"
                }
            ],
            "versionInfo": {
                "agentVersion": "1.8.2",
                "agentHash": "64205ac",
                "dockerVersion": "DockerVersion: 1.9.1"
            }
        }
    ]
}
`},
}
