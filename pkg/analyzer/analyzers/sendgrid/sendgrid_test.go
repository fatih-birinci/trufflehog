package sendgrid

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/trufflesecurity/trufflehog/v3/pkg/analyzer/analyzers"
	"github.com/trufflesecurity/trufflehog/v3/pkg/analyzer/config"
	"github.com/trufflesecurity/trufflehog/v3/pkg/common"
	"github.com/trufflesecurity/trufflehog/v3/pkg/context"
)

func TestAnalyzer_Analyze(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	testSecrets, err := common.GetSecret(ctx, "trufflehog-testing", "detectors5")
	if err != nil {
		t.Fatalf("could not get test secrets from GCP: %s", err)
	}

	tests := []struct {
		name    string
		key     string
		want    string // JSON string
		wantErr bool
	}{
		{
			name: "Valid Sendgrid key",
			key:  testSecrets.MustGetField("SENDGRID"),
			want: `{
				"AnalyzerType":14,
				"Bindings":[
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"api_keys.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"api_keys.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"api_keys.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"api_keys.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"alerts.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"alerts.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"alerts.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"alerts.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"categories.stats.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"categories.stats.sums.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"categories.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"categories.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"categories.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"categories.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"design_library.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"design_library.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"design_library.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"design_library.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"messages.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"email_testing.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"email_testing.write",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"ips.pools.ips.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.webhooks.parse.settings.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.webhooks.parse.settings.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.webhooks.parse.settings.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.webhooks.parse.settings.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail.send",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.address_whitelist.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.address_whitelist.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.address_whitelist.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.address_whitelist.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.bcc.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.bcc.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.bcc.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.bcc.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.bounce_purge.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.bounce_purge.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.bounce_purge.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.bounce_purge.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.webhooks.event.settings.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.webhooks.event.settings.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.webhooks.event.settings.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.webhooks.event.settings.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.webhooks.event.test.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.webhooks.event.test.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.webhooks.event.test.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.webhooks.event.test.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.footer.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.footer.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.footer.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.footer.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.forward_bounce.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.forward_bounce.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.forward_bounce.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.forward_bounce.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.forward_spam.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.forward_spam.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.forward_spam.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.forward_spam.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.template.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.template.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.template.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.template.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.plain_content.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.plain_content.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.plain_content.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.plain_content.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.spam_check.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.spam_check.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.spam_check.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.spam_check.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"marketing.automation.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"marketing.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"partner_settings.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"partner_settings.new_relic.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"partner_settings.new_relic.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"partner_settings.new_relic.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"partner_settings.new_relic.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"partner_settings.sendwithus.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"partner_settings.sendwithus.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"partner_settings.sendwithus.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"partner_settings.sendwithus.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"recipients.erasejob.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"recipients.erasejob.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"access_settings.whitelist.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"access_settings.whitelist.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"access_settings.whitelist.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"access_settings.whitelist.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"access_settings.activity.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"whitelabel.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"whitelabel.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"whitelabel.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"whitelabel.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"browsers.stats.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"devices.stats.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"clients.stats.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"clients.phone.stats.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"clients.tablet.stats.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"clients.webmail.stats.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"clients.desktop.stats.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"geo.stats.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"stats.global.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mailbox_providers.stats.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.webhooks.parse.stats.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"stats.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"subusers.stats.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"subusers.stats.sums.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"subusers.stats.monthly.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"suppression.bounces.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"suppression.bounces.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"suppression.bounces.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"suppression.bounces.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"suppression.blocks.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"suppression.blocks.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"suppression.blocks.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"suppression.blocks.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"suppression.invalid_emails.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"suppression.invalid_emails.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"suppression.invalid_emails.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"suppression.invalid_emails.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"suppression.spam_reports.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"suppression.spam_reports.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"suppression.spam_reports.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"suppression.spam_reports.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"suppression.unsubscribes.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"suppression.unsubscribes.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"suppression.unsubscribes.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"suppression.unsubscribes.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"suppression.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"suppression.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"suppression.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"suppression.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"asm.groups.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"asm.groups.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"asm.groups.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"asm.groups.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"asm.groups.suppressions.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"asm.groups.suppressions.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"asm.groups.suppressions.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"asm.groups.suppressions.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"templates.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"templates.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"templates.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"templates.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"templates.versions.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"templates.versions.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"templates.versions.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"templates.versions.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"templates.versions.activate.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"templates.versions.activate.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"templates.versions.activate.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"templates.versions.activate.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"tracking_settings.click.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"tracking_settings.click.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"tracking_settings.click.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"tracking_settings.click.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"tracking_settings.google_analytics.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"tracking_settings.google_analytics.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"tracking_settings.google_analytics.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"tracking_settings.google_analytics.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"tracking_settings.open.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"tracking_settings.open.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"tracking_settings.open.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"tracking_settings.open.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"tracking_settings.subscription.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"tracking_settings.subscription.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"tracking_settings.subscription.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"tracking_settings.subscription.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.settings.enforced_tls.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.settings.enforced_tls.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.timezone.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.timezone.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.timezone.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.timezone.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"asm.suppressions.global.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"asm.suppressions.global.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"asm.suppressions.global.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"asm.suppressions.global.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"credentials.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"credentials.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"credentials.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"credentials.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"mail_settings.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"signup.trigger_confirmation",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"teammates.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"teammates.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"teammates.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"teammates.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"tracking_settings.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"ui.confirm_email",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"ui.provision",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"ui.signup_complete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.account.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.credits.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.email.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.email.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.email.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.email.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.multifactor_authentication.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.multifactor_authentication.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.multifactor_authentication.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.multifactor_authentication.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.password.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.password.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.password.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.password.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.profile.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.profile.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.profile.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.profile.delete",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.username.create",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.username.read",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.username.update",
						 "Parent":null
					  }
				   },
				   {
					  "Resource":{
						 "Name":"Sendgrid Key",
						 "FullyQualifiedName":"Sendgrid Key",
						 "Type":"key",
						 "Metadata":{
							"2fa_required":true,
							"key_type":"full access"
						 },
						 "Parent":null
					  },
					  "Permission":{
						 "Value":"user.username.delete",
						 "Parent":null
					  }
				   }
				],
				"UnboundedResources":null,
				"Metadata":null
			 }`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Analyzer{Cfg: &config.Config{}}
			got, err := a.Analyze(ctx, map[string]string{"key": tt.key})
			if (err != nil) != tt.wantErr {
				t.Errorf("Analyzer.Analyze() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Marshal the actual result to JSON
			gotJSON, err := json.Marshal(got)
			if err != nil {
				t.Fatalf("could not marshal got to JSON: %s", err)
			}

			// Parse the expected JSON string
			var wantObj analyzers.AnalyzerResult
			if err := json.Unmarshal([]byte(tt.want), &wantObj); err != nil {
				t.Fatalf("could not unmarshal want JSON string: %s", err)
			}

			// Marshal the expected result to JSON (to normalize)
			wantJSON, err := json.Marshal(wantObj)
			if err != nil {
				t.Fatalf("could not marshal want to JSON: %s", err)
			}

			// Compare the JSON strings
			if string(gotJSON) != string(wantJSON) {
				// Pretty-print both JSON strings for easier comparison
				var gotIndented, wantIndented []byte
				gotIndented, err = json.MarshalIndent(got, "", " ")
				if err != nil {
					t.Fatalf("could not marshal got to indented JSON: %s", err)
				}
				wantIndented, err = json.MarshalIndent(wantObj, "", " ")
				if err != nil {
					t.Fatalf("could not marshal want to indented JSON: %s", err)
				}
				t.Errorf("Analyzer.Analyze() = %s, want %s", gotIndented, wantIndented)
			}
		})
	}
}
