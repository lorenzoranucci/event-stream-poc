<?php

declare(strict_types=1);

namespace App\MerchantReview\Domain;

final class MerchantReviewUUID
{
    private function __construct(
        private string $value,
    ) {}

    public static function createFromString(
        string $value,
    ): self {
        return new self($value);
    }

    public function __toString(): string
    {
        return $this->value;
    }
}
