import 'account.dart';

enum Provider {
  lifx('lifx', 'LIFX'),
  hue('hue', 'Philips Hue');

  const Provider(this.value, this.displayName);

  final String value;
  final String displayName;

  static Provider fromString(String value) {
    return Provider.values.firstWhere(
      (provider) => provider.value == value,
      orElse: () => throw ArgumentError('Invalid provider: $value'),
    );
  }
}

class ConnectProviderRequest {
  final String provider;
  final String token;

  ConnectProviderRequest({
    required this.provider,
    required this.token,
  });

  Map<String, dynamic> toJson() {
    return {
      'provider': provider,
      'token': token,
    };
  }
}

class ConnectProviderResponse {
  final String id;
  final String provider;
  final String providerAccountId;
  final Map<String, dynamic>? metadata;
  final DateTime createdAt;

  ConnectProviderResponse({
    required this.id,
    required this.provider,
    required this.providerAccountId,
    this.metadata,
    required this.createdAt,
  });

  factory ConnectProviderResponse.fromJson(Map<String, dynamic> json) {
    return ConnectProviderResponse(
      id: json['id'] as String,
      provider: json['provider'] as String,
      providerAccountId: json['provider_account_id'] as String,
      metadata: json['metadata'] as Map<String, dynamic>?,
      createdAt: DateTime.parse(json['created_at'] as String),
    );
  }
}

class ListAccountsResponse {
  final List<Account> accounts;

  ListAccountsResponse({required this.accounts});

  factory ListAccountsResponse.fromJson(Map<String, dynamic> json) {
    final accountsList = json['accounts'] as List<dynamic>;
    return ListAccountsResponse(
      accounts: accountsList
          .map((account) => Account.fromJson(account as Map<String, dynamic>))
          .toList(),
    );
  }
}
