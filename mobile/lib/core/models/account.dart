class Account {
  final String id;
  final String provider;
  final String providerAccountId;
  final Map<String, dynamic>? metadata;
  final DateTime createdAt;

  Account({
    required this.id,
    required this.provider,
    required this.providerAccountId,
    this.metadata,
    required this.createdAt,
  });

  factory Account.fromJson(Map<String, dynamic> json) {
    return Account(
      id: json['id'] as String,
      provider: json['provider'] as String,
      providerAccountId: json['provider_account_id'] as String,
      metadata: json['metadata'] as Map<String, dynamic>?,
      createdAt: DateTime.parse(json['created_at'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'provider': provider,
      'provider_account_id': providerAccountId,
      'metadata': metadata,
      'created_at': createdAt.toIso8601String(),
    };
  }

  Account copyWith({
    String? id,
    String? provider,
    String? providerAccountId,
    Map<String, dynamic>? metadata,
    DateTime? createdAt,
  }) {
    return Account(
      id: id ?? this.id,
      provider: provider ?? this.provider,
      providerAccountId: providerAccountId ?? this.providerAccountId,
      metadata: metadata ?? this.metadata,
      createdAt: createdAt ?? this.createdAt,
    );
  }

  String get displayName {
    // Try to get a label from metadata first
    if (metadata != null && metadata!.containsKey('label')) {
      return metadata!['label'] as String;
    }
    // Fallback to provider name
    return provider.toUpperCase();
  }
}
